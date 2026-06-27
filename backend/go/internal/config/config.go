package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret            string
	JWTExpirationMinutes int
	HasuraURL            string
	HasuraAdminSecret    string
	Port                 string
	ChapaSecretKey       string
	RecipeBaseURL		 string
	SMTPEmail            string
	SMTPPassword         string
	CloudinaryURL        string
}

func LoadConfig() Config {
	// Load the .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found.")
	}

	return Config{
		JWTSecret:            getRequiredEnv("JWT_SECRET"),
		JWTExpirationMinutes: getRequiredEnvAsInt("JWT_EXPIRATION_MINUTES"),
		HasuraURL:            getRequiredEnv("HASURA_URL"),
		HasuraAdminSecret:    getRequiredEnv("HASURA_ADMIN_SECRET"),
		Port:                 getRequiredEnv("PORT"),
		ChapaSecretKey:       getRequiredEnv("CHAPA_SECRET_KEY"),
		RecipeBaseURL:		  getRequiredEnv("RECIPE_BASE_URL"),
		SMTPEmail:            getRequiredEnv("SMTP_EMAIL"),
		SMTPPassword:         getRequiredEnv("SMTP_PASSWORD"),
		CloudinaryURL:        getRequiredEnv("CLOUDINARY_URL"),
	}
}

// getRequiredEnv fetches a string or crashes the app if missing
func getRequiredEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("FATAL: Environment variable %s is missing or empty", key)
	}
	return value
}

// getRequiredEnvAsInt fetches an integer or crashes the app if missing/invalid
func getRequiredEnvAsInt(key string) int {
	valueStr := getRequiredEnv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("FATAL: Environment variable %s must be a valid integer, got '%s'", key, valueStr)
	}
	return value
}
