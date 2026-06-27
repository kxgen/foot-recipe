package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/internal/config"
	"server/internal/dto/payment"
	"time"

	graphql "github.com/hasura/go-graphql-client"
)

type Numeric float64

func (n Numeric) GetGraphQLType() string {
	return "numeric"
}

type PaymentService struct {
	hasura *graphql.Client
	config config.Config
}

func NewPaymentService(hasura *graphql.Client, cfg config.Config) *PaymentService {
	return &PaymentService{
		hasura: hasura,
		config: cfg,
	}
}

func (s *PaymentService) InitializeChapaPayment(ctx context.Context, recipeID, userID int, price float64, userEmail, userName, txRef string) (string, error) {
	// 1. Call Chapa to initialize payment
	chapaReq := payment.ChapaInitializeRequest{
		Amount:        fmt.Sprintf("%.2f", price),
		Currency:      "ETB",
		Email:         userEmail,
		FirstName:     userName,
		LastName:      "",
		TxRef:         txRef,
		ReturnURL:     fmt.Sprintf("http://localhost:8000/payments/chapa/verify?tx_ref=%s&recipe_id=%d", txRef, recipeID),
		Customization: map[string]string{
		    "title":       "Recipe Purchase",
		    "description": fmt.Sprintf("Payment for Recipe %d", recipeID),
		},
	}

	reqBody, err := json.Marshal(chapaReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chapa request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.chapa.co/v1/transaction/initialize", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.config.ChapaSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call chapa initialize: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("chapa initialize failed with status %d: %s", resp.StatusCode, string(body))
	}

	var chapaResp payment.ChapaInitializeResponse
	if err := json.NewDecoder(resp.Body).Decode(&chapaResp); err != nil {
		return "", fmt.Errorf("failed to decode chapa response: %w", err)
	}

	if chapaResp.Status != "success" {
		return "", fmt.Errorf("chapa returned unsuccessful status: %s", chapaResp.Message)
	}

	// 2. Insert pending purchase record in Hasura
	var mut struct {
		InsertRecipePurchasesOne struct {
			ID int `graphql:"id"`
		} `graphql:"insert_recipe_purchases_one(object: {recipe_id: $recipe_id, user_id: $user_id, amount_paid: $amount, status: \"pending\", payment_provider: \"chapa\", provider_transaction_id: $tx_ref}, on_conflict: {constraint: recipe_purchases_user_id_recipe_id_key, update_columns: [status, amount_paid, provider_transaction_id, payment_provider]})"`
	}

	vars := map[string]interface{}{
			"recipe_id": graphql.Int(recipeID),
			"user_id":   graphql.Int(userID),
			"amount":    Numeric(price),
			"tx_ref":    graphql.String(txRef),
	}

	if err := s.hasura.Mutate(ctx, &mut, vars); err != nil {
		return "", fmt.Errorf("failed to insert pending purchase: %w", err)
	}

	return chapaResp.Data.CheckoutURL, nil
}

func (s *PaymentService) VerifyChapaPayment(ctx context.Context, txRef string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.chapa.co/v1/transaction/verify/"+txRef, nil)
	if err != nil {
		return fmt.Errorf("failed to create verify request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.config.ChapaSecretKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call chapa verify: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("chapa verify failed with status %d: %s", resp.StatusCode, string(body))
	}

	var verifyResp payment.ChapaVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return fmt.Errorf("failed to decode chapa verify response: %w", err)
	}

	if verifyResp.Status != "success" || verifyResp.Data.Status != "success" {
		return fmt.Errorf("payment not successful according to chapa: %s", verifyResp.Message)
	}

	var mut struct {
		UpdateRecipePurchases struct {
			AffectedRows int `graphql:"affected_rows"`
		} `graphql:"update_recipe_purchases(where: {provider_transaction_id: {_eq: $tx_ref}}, _set: {status: \"completed\"})"`
	}

	vars := map[string]interface{}{
		"tx_ref": graphql.String(txRef),
	}

	if err := s.hasura.Mutate(ctx, &mut, vars); err != nil {
		return fmt.Errorf("failed to update purchase status in hasura: %w", err)
	}

	if mut.UpdateRecipePurchases.AffectedRows == 0 {
		return fmt.Errorf("no pending purchase found for tx_ref: %s", txRef)
	}

	return nil
}
