package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	graphql "github.com/hasura/go-graphql-client"
)

type RecipeService struct {
	hasura *graphql.Client
	cld    *cloudinary.Cloudinary
}

func NewRecipeService(hasura *graphql.Client, cld *cloudinary.Cloudinary) *RecipeService {
	return &RecipeService{hasura: hasura, cld: cld}
}

// ProcessAndSaveRecipeImage uploads a base64 image string directly to Cloudinary
// under the shareplate/recipes folder and returns the public secure URL immediately.
func (s *RecipeService) ProcessAndSaveRecipeImage(imageName string, base64Image string, recipeId int, tempFolder string, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Direct upload with unique ID to avoid collisions
	publicID := fmt.Sprintf("img_%d_%d", time.Now().UnixNano(), rand.Intn(10000))

	resp, err := s.cld.Upload.Upload(ctx, base64Image, uploader.UploadParams{
		Folder:         "shareplate/recipes",
		PublicID:       publicID,
		Overwrite:      boolPtr(true),
		Invalidate:     boolPtr(true),
		UniqueFilename: boolPtr(false),
	})
	if err != nil {
		log.Printf("Cloudinary upload error: %v", err)
		return "", fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	return resp.SecureURL, nil
}

