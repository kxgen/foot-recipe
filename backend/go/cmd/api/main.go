package main

import (
	"log"
	"net/http"
	"server/internal/config"
	"server/internal/handler"
	"server/internal/hasura"
	"server/internal/service"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Hasura client
	hasuraClient := hasura.NewClient(cfg)

	// Initialize Cloudinary client
	cld, err := cloudinary.NewFromURL(cfg.CloudinaryURL)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize Cloudinary client: %v", err)
	}

	// Initialize services
	authService := service.NewAuthService(hasuraClient, cfg)
	userService := service.NewUserService(hasuraClient, cld)
	recipeService := service.NewRecipeService(hasuraClient, cld)
	paymentService := service.NewPaymentService(hasuraClient, cfg)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	recipeHandler := handler.NewRecipeHandler(recipeService)
	paymentHandler := handler.NewPaymentHandler(paymentService, hasuraClient, cfg)

	// 1. Initialize ServeMux router
	mux := http.NewServeMux()

	// 2. Register static file serving on the isolated mux
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// 3. Define routes on the isolated mux instead of the global http object
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/users/upload-avatar", userHandler.UploadAvatarHandler)
	mux.HandleFunc("/users/update-password", userHandler.UpdatePasswordHandler)

	mux.HandleFunc("/recipes/upload-image", recipeHandler.UploadRecipeImageHandler)

	mux.HandleFunc("/payments/chapa/initialize", paymentHandler.InitializeChapaPaymentHandler)
	mux.HandleFunc("/payments/chapa/verify", paymentHandler.VerifyChapaPaymentHandler)
	mux.HandleFunc("/payments/send-email", paymentHandler.SendConfirmationEmailHandler)

	// 4. Configure custom production-ready server settings
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,   // Max duration for reading the request
		WriteTimeout: 10 * time.Second,  // Max duration for writing the response
		IdleTimeout:  120 * time.Second, // Max time to keep idle connections alive
	}

	log.Printf("Server running on :%s", cfg.Port)

	// 5. Start the server using the configured struct
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("FATAL: Server failed to start: %v", err)
	}
}
