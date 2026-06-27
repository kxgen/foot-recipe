package handler

import (
	"encoding/json"
	"net/http"
	"server/internal/dto/recipe"
	"server/internal/service"
)

type RecipeHandler struct {
	recipeService *service.RecipeService
}

func NewRecipeHandler(recipeService *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{
		recipeService: recipeService,
	}
}

// Webhook Handler for uploading recipe images
func (h *RecipeHandler) UploadRecipeImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Enforce POST method
	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 2. Decode Hasura Action payload
	var payload struct {
		Action struct {
			Name string `json:"name"`
		} `json:"action"`
		Input recipe.UploadRecipeImageRequest `json:"input"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	input := payload.Input

	// 3. Validate required fields
	if input.ImageName == "" || input.Base64Image == "" {
		h.respondWithError(w, http.StatusBadRequest, "Missing required fields: imageName and base64Image are required")
		return
	}

	// 4. Save image and return its URL
	imageUrl, err := h.recipeService.ProcessAndSaveRecipeImage(input.ImageName, input.Base64Image, input.RecipeId, input.TempFolder, input.Filename)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 5. Send success response
	w.WriteHeader(http.StatusOK)
	res := recipe.UploadRecipeImageResponse{
		Success:  true,
		ImageUrl: imageUrl,
	}
	json.NewEncoder(w).Encode(res)
}



// respondWithError helper writes a JSON error response
func (h *RecipeHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
