package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/internal/dto/user"
	"server/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

// Constructor for User Handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Webhook Handler func for Uploading Avatar 
func (h *UserHandler) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Enforce POST method
	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 2. Decode incoming Hasura -> to Anonymous Struct payload{action{name}} 
	var payload struct {
		// Hasura Action Name -> uploadAvatar
		Action struct {
			Name string `json:"name"`
		} `json:"action"`
		// Sent Input
		Input user.UploadAvatarRequest `json:"input"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	input := payload.Input

	// 3. Validate required fields
	if input.UserId <= 0 || input.Slug == "" || input.ImageName == "" || input.Base64Image == "" {
		h.respondWithError(w, http.StatusBadRequest, "Missing required fields: userId, slug, imageName, and base64Image are required")
		return
	}

	// 4. Process image and save to storage & DB
	avatarUrl, err := h.userService.ProcessAndSaveAvatar(input.UserId, input.Slug, input.ImageName, input.Base64Image)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 5. Send success response
	w.WriteHeader(http.StatusOK)
	res := user.UploadAvatarResponse{
		Success:   true,
		AvatarUrl: avatarUrl,
	}
	json.NewEncoder(w).Encode(res)
}

// UpdatePasswordHandler handles the Hasura action for updating user password
func (h *UserHandler) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var payload struct {
		Action struct {
			Name string `json:"name"`
		} `json:"action"`
		Input            user.UpdatePasswordRequest `json:"input"`
		SessionVariables map[string]string          `json:"session_variables"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Get UserID from session variables
	userIdStr, ok := payload.SessionVariables["x-hasura-user-id"]
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "User ID not found in session variables")
		return
	}

	// Convert userIdStr to int
	var userId int
	if _, err := fmt.Sscanf(userIdStr, "%d", &userId); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid User ID in session variables")
		return
	}

	if payload.Input.NewPassword == "" {
		h.respondWithError(w, http.StatusBadRequest, "New password is required")
		return
	}

	err := h.userService.UpdatePassword(r.Context(), userId, payload.Input.NewPassword)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	res := user.UpdatePasswordResponse{
		Success: true,
		Message: "Password updated successfully",
	}
	json.NewEncoder(w).Encode(res)
}

// Helper Function to write JSON error response
func (h *UserHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
