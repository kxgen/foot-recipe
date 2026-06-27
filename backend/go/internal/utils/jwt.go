package utils

import (
	"strconv"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type JWTInput struct {
	UserID int
}

type JWTConfig struct {
	SecretKey           string
	ExpirationMinutes   int
}

// all signed in users
var ROLE = "user"

func GenerateJWT(input JWTInput, config JWTConfig) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": strconv.Itoa(input.UserID),
		"iat": now.Unix(),
		"exp": now.Add(time.Minute * time.Duration(config.ExpirationMinutes)).Unix(),
		"https://hasura.io/jwt/claims": map[string]interface{}{
			"x-hasura-default-role":  ROLE,
			"x-hasura-allowed-roles": []string{ROLE},
			"x-hasura-user-id":       strconv.Itoa(input.UserID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
