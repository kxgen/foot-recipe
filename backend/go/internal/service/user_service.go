package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"server/internal/utils"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	graphql "github.com/hasura/go-graphql-client"
)

type UserService struct {
	hasura *graphql.Client
	cld    *cloudinary.Cloudinary
}

func NewUserService(hasura *graphql.Client, cld *cloudinary.Cloudinary) *UserService {
	return &UserService{hasura: hasura, cld: cld}
}

func boolPtr(b bool) *bool {
	return &b
}

// ProcessAndSaveAvatar uploads a base64 avatar image string directly to Cloudinary
// under the shareplate/avatars folder and updates the user's avatar_url in Hasura.
func (s *UserService) ProcessAndSaveAvatar(userId int, slug string, imageName string, base64Image string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Upload to Cloudinary under avatars folder
	publicID := fmt.Sprintf("avatar-%s", slug)
	resp, err := s.cld.Upload.Upload(ctx, base64Image, uploader.UploadParams{
		Folder:         "shareplate/avatars",
		PublicID:       publicID,
		Overwrite:      boolPtr(true),
		Invalidate:     boolPtr(true),
		UniqueFilename: boolPtr(false),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar to Cloudinary: %w", err)
	}

	avatarUrl := resp.SecureURL

	// 2. Update user's avatar_url in database
	var mutation struct {
		UpdateUsersByPk struct {
			ID int `graphql:"id"`
		} `graphql:"update_users_by_pk(pk_columns: {id: $id}, _set: {avatar_url: $avatarUrl})"`
	}

	variables := map[string]interface{}{
		"id":        graphql.Int(userId),
		"avatarUrl": graphql.String(avatarUrl),
	}

	if err := s.hasura.Mutate(ctx, &mutation, variables); err != nil {
		return "", fmt.Errorf("hasura database update failed: %w", err)
	}

	return avatarUrl, nil
}

// UpdatePassword hashes the new password and updates the user's password in the database.
func (s *UserService) UpdatePassword(ctx context.Context, userId int, newPassword string) error {
	// 1. Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 2. Update the user's password in Hasura
	var mutation struct {
		UpdateUsersByPk struct {
			ID int `graphql:"id"`
		} `graphql:"update_users_by_pk(pk_columns: {id: $id}, _set: {password: $password})"`
	}

	variables := map[string]interface{}{
		"id":       graphql.Int(userId),
		"password": graphql.String(hashedPassword),
	}

	if err := s.hasura.Mutate(ctx, &mutation, variables); err != nil {
		return fmt.Errorf("failed to update password in database: %w", err)
	}

	return nil
}

// helper function to decode base64 strings
func decodeBase64(base64Str string) ([]byte, error) {
	parts := strings.Split(base64Str, ",")
	var rawBase64 string
	if len(parts) > 1 {
		rawBase64 = parts[1]
	} else {
		rawBase64 = base64Str
	}

	// Remove any white spaces
	rawBase64 = strings.TrimSpace(rawBase64)

	data, err := base64.StdEncoding.DecodeString(rawBase64)
	if err != nil {
		return nil, err
	}
	return data, nil
}
