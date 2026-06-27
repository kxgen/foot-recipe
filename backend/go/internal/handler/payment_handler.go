package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"server/internal/config"
	"server/internal/service"
	"time"

	graphql "github.com/hasura/go-graphql-client"
	"gopkg.in/gomail.v2"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
	hasura         *graphql.Client
	config         config.Config
}

func NewPaymentHandler(paymentService *service.PaymentService, hasura *graphql.Client, cfg config.Config) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		hasura:         hasura,
		config:         cfg,
	}
}

// InitializeChapaPaymentHandler handles the Hasura Action webhook for initializing a Chapa payment
func (h *PaymentHandler) InitializeChapaPaymentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Enforce POST method
	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 2. Decode the incoming Hasura Action payload {action{} & input {}}
	var payload struct {
		Action struct {
			Name string `json:"name"`
		} `json:"action"`
		Input struct {
			RecipeID int `json:"recipeId"`
			UserID   int `json:"userId"`
		} `json:"input"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	req := payload.Input

	if req.RecipeID <= 0 || req.UserID <= 0 {
		h.respondWithError(w, http.StatusBadRequest, "Missing recipeId or userId in input")
		return
	}

	// Fetch recipe details and user details from Hasura using Request Context
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var query struct {
		RecipesByPk struct {
			Price   float64 `graphql:"price"`
			Title   string  `graphql:"title"`
			User_ID int     `graphql:"user_id"`
		} `graphql:"recipes_by_pk(id: $recipe_id)"`
		UsersByPk struct {
			Username string `graphql:"username"`
			Email    string `graphql:"email"`
		} `graphql:"users_by_pk(id: $user_id)"`
	}

	vars := map[string]interface{}{
		"recipe_id": graphql.Int(req.RecipeID),
		"user_id":   graphql.Int(req.UserID),
	}

	if err := h.hasura.Query(ctx, &query, vars); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to query Hasura: "+err.Error())
		return
	}

	if query.RecipesByPk.Title == "" {
		h.respondWithError(w, http.StatusNotFound, "Recipe not found")
		return
	}

	if query.RecipesByPk.User_ID == req.UserID {
		h.respondWithError(w, http.StatusBadRequest, "Author cannot purchase their own recipe")
		return
	}

	if query.UsersByPk.Email == "" {
		h.respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Generate tx_ref
	txRef := fmt.Sprintf("chapa-%d-%d-%d", req.RecipeID, req.UserID, time.Now().Unix())

	// Initialize payment
	checkoutUrl, err := h.paymentService.InitializeChapaPayment(ctx, req.RecipeID, req.UserID, query.RecipesByPk.Price, query.UsersByPk.Email, query.UsersByPk.Username, txRef)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to initialize payment: "+err.Error())
		return
	}

	res := map[string]string{
		"checkoutUrl": checkoutUrl,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *PaymentHandler) VerifyChapaPaymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	txRef := r.URL.Query().Get("tx_ref")
	if txRef == "" {
		txRef = r.URL.Query().Get("trx_ref")
	}
	recipeIdStr := r.URL.Query().Get("recipe_id")

	log.Printf("VerifyChapaPayment: tx_ref=%s, recipe_id=%s", txRef, recipeIdStr)

	// If recipe_id is missing, try to extract it from tx_ref (chapa-RECIPEID-USERID-TIMESTAMP)
	if recipeIdStr == "" && txRef != "" {
		parts := strings.Split(txRef, "-")
		if len(parts) >= 2 {
			recipeIdStr = parts[1]
			log.Printf("Extracted recipe_id %s from tx_ref %s", recipeIdStr, txRef)
		}
	}

	// Sanitize recipe_id to prevent injection vulnerabilities
	if _, err := strconv.Atoi(recipeIdStr); err != nil || txRef == "" {
		log.Printf("Verification failed: missing or invalid params. tx_ref=%s, recipe_id=%s", txRef, recipeIdStr)
		http.Error(w, "Invalid or missing tx_ref or recipe_id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := h.paymentService.VerifyChapaPayment(ctx, txRef)
	if err != nil {
		fmt.Printf("Payment verification failed for %s: %v\n", txRef, err)
		http.Redirect(w, r, h.config.RecipeBaseURL+recipeIdStr+"?payment=failed", http.StatusFound)
		return
	}

	log.Printf("Payment verification successful for %s. Redirecting to recipe %s", txRef, recipeIdStr)

	// Success, redirect to frontend
	http.Redirect(w, r, h.config.RecipeBaseURL+recipeIdStr+"?payment=success", http.StatusFound)
}

func (h *PaymentHandler) SendConfirmationEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var payload struct {
		Event struct {
			Data struct {
				New struct {
					UserID int `json:"user_id"`
				} `json:"new"`
			} `json:"data"`
		} `json:"event"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	userID := payload.Event.Data.New.UserID
	if userID <= 0 {
		h.respondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	// Query Hasura for user email
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var query struct {
		UsersByPk struct {
			Email string `graphql:"email"`
		} `graphql:"users_by_pk(id: $user_id)"`
	}
	vars := map[string]interface{}{
		"user_id": graphql.Int(userID),
	}

	if err := h.hasura.Query(ctx, &query, vars); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to query user: "+err.Error())
		return
	}

	userEmail := query.UsersByPk.Email
	if userEmail == "" {
		h.respondWithError(w, http.StatusNotFound, "User email not found")
		return
	}

	// Send email using gomail
	m := gomail.NewMessage()
	m.SetHeader("From", h.config.SMTPEmail)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Payment Confirmation")
	m.SetBody("text/plain", "Thank you for your purchase! Your payment has been confirmed.")

	log.Printf("Email confirmation triggered for user email: %s\n", userEmail)

	d := gomail.NewDialer("smtp.gmail.com", 587, h.config.SMTPEmail, h.config.SMTPPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v\n", userEmail, err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to send email")
		return
	}

	log.Printf("Email confirmation successfully sent to: %s\n", userEmail)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *PaymentHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
