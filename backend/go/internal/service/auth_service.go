package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"server/internal/config"
	"server/internal/dto/auth"
	"server/internal/utils"
	"strings"

	graphql "github.com/hasura/go-graphql-client"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserExists         = errors.New("username or email already exists")
)

type AuthService struct {
	hasura *graphql.Client
	cfg    config.Config
}

func NewAuthService(hasura *graphql.Client, cfg config.Config) *AuthService {
	return &AuthService{
		hasura: hasura,
		cfg:    cfg,
	}
}

func (s *AuthService) Login(ctx context.Context, input auth.LoginRequest) (auth.AuthResponse, error) {
	var query struct {
		Users []struct {
			ID       int    `graphql:"id"`
			Password string `graphql:"password"`
		} `graphql:"users(where: {username: {_eq: $username}}, limit: 1)"`
	}

	variables := map[string]interface{}{
		"username": graphql.String(input.Username),
	}

	if err := s.hasura.Query(ctx, &query, variables); err != nil {
		return auth.AuthResponse{}, fmt.Errorf("failed to query database: %w", err)
	}

	if len(query.Users) == 0 {
		return auth.AuthResponse{}, ErrInvalidCredentials
	}

	user := query.Users[0]
	
	// bcrypt comparison for check validation
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return auth.AuthResponse{}, ErrInvalidCredentials
	}

	token, err := utils.GenerateJWT(utils.JWTInput{UserID: user.ID}, utils.JWTConfig{
		SecretKey:         s.cfg.JWTSecret,
		ExpirationMinutes: s.cfg.JWTExpirationMinutes,
	})
	if err != nil {
		return auth.AuthResponse{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return auth.AuthResponse{ID: user.ID, Token: token}, nil
}

func (s *AuthService) Register(ctx context.Context, input auth.RegisterRequest) (auth.AuthResponse, error) {
	// 1. Check if the user already exists
	var checkQuery struct {
		Users []struct {
			ID int `graphql:"id"`
		} `graphql:"users(where: {_or: [{username: {_eq: $username}}, {email: {_eq: $email}}]}, limit: 1)"`
	}

	variables := map[string]interface{}{
		"username": graphql.String(input.Username),
		"email":    graphql.String(input.Email),
	}

	if err := s.hasura.Query(ctx, &checkQuery, variables); err != nil {
		return auth.AuthResponse{}, fmt.Errorf("failed to check existing user: %w", err)
	}

	if len(checkQuery.Users) > 0 {
		return auth.AuthResponse{}, ErrUserExists
	}

	// Hash the incoming plain text password before pushing to Hasura
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return auth.AuthResponse{}, fmt.Errorf("security error: %w", err)
	}

	// Generate Slug
	slug := s.slugify(input.Username)

	// 2. Define the GraphQL Mutation to insert the new user
	var mutation struct {
		InsertUsersOne struct {
			ID int `graphql:"id"`
		} `graphql:"insert_users_one(object: {username: $username, password: $password, email: $email, slug: $slug})"`
	}

	mutationVariables := map[string]interface{}{
		"username": graphql.String(input.Username),
		"password": graphql.String(hashedPassword),
		"email":    graphql.String(input.Email),
		"slug":     graphql.String(slug),
	}

	// 3. Execute the mutation on Hasura
	if err := s.hasura.Mutate(ctx, &mutation, mutationVariables); err != nil {
		return auth.AuthResponse{}, fmt.Errorf("failed to insert user: %w", err)
	}

	newUserID := mutation.InsertUsersOne.ID

	// 4. Generate JWT for the freshly registered user
	token, err := utils.GenerateJWT(utils.JWTInput{UserID: newUserID}, utils.JWTConfig{
		SecretKey:         s.cfg.JWTSecret,
		ExpirationMinutes: s.cfg.JWTExpirationMinutes,
	})
	if err != nil {
		return auth.AuthResponse{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return auth.AuthResponse{ID: newUserID, Token: token}, nil
}

func (s *AuthService) slugify(str string) string {
	str = strings.ToLower(str)
	str = strings.TrimSpace(str)
	// Replace spaces with -
	str = strings.ReplaceAll(str, " ", "-")
	// Remove non-alphanumeric (except -)
	reg, _ := regexp.Compile("[^a-z0-9-]+")
	str = reg.ReplaceAllString(str, "")
	// Remove multiple dashes
	reg2, _ := regexp.Compile("-+")
	str = reg2.ReplaceAllString(str, "-")
	return str
}