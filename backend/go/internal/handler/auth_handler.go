package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/internal/dto/auth"
	"server/internal/service"
)

// dependency
type AuthHandler struct {
	authService *service.AuthService
}

// Constructor for an Instance of Auth Handler
func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

// Login Handler
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Enforce POST Method
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, HasuraError{Message: "Method not allowed"})
		return
	}

	// Decode Incoming Request to LoginRequest DTO
	var body struct {
		Input auth.LoginRequest `json:"input"`
	}

	// Incoming r.Body decoded --to-> var body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, HasuraError{Message: "Invalid request payload"})
		return
	}

	// Pass the credentials down to service layer
	res, err := h.authService.Login(r.Context(), body.Input)
	if err != nil {
		// If credentials are bad, return 401 Unauthorized
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeJSON(w, http.StatusUnauthorized, HasuraError{Message: err.Error()})
			return
		}
		// Catch-all for unexpected database/internal issues
		writeJSON(w, http.StatusInternalServerError, HasuraError{Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Enforce POST Method
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, HasuraError{Message: "Method not allowed"})
		return
	}

	// Decode Incoming Request to RegisterRequest DTO
	var body struct {
		Input auth.RegisterRequest `json:"input"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, HasuraError{Message: "Invalid request payload"})
		return
	}

	// Pass registration payload to Register method
	res, err := h.authService.Register(r.Context(), body.Input)
	if err != nil {
		// If the user already exists, return 409 Conflict
		if errors.Is(err, service.ErrUserExists) {
			writeJSON(w, http.StatusConflict, HasuraError{Message: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, HasuraError{Message: err.Error()})
		return
	}
	
	writeJSON(w, http.StatusOK, res)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Hasura's custom error format
type HasuraError struct {
	Message string `json:"message"`
}