package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/lcleigh/page-hoppers-backend/handlers"
    "github.com/lcleigh/page-hoppers-backend/models"
	gorillahandlers "github.com/gorilla/handlers"
)

func main() {
	// Load environment variables
	if os.Getenv("IN_DOCKER") != "true" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found (this is fine if running in Docker)")
		}
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&models.User{}, &models.ReadingLog{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize router
	r := mux.NewRouter()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, []byte(os.Getenv("JWT_SECRET")))
	readingLogHandler := handlers.NewReadingLogHandler(db)

	log.Println("Initializing routes...")

	// Public routes
	r.HandleFunc("/api/auth/parent/login", authHandler.ParentLogin).Methods("POST")
	r.HandleFunc("/api/auth/parent/register", authHandler.ParentRegister).Methods("POST")
	r.HandleFunc("/api/auth/child/login", authHandler.ChildLogin).Methods("POST")

	// Protected routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(authMiddleware)

	// Parent-only routes
	api.HandleFunc("/children", authHandler.GetChildren).Methods("GET")
	api.HandleFunc("/children", authHandler.CreateChild).Methods("POST")

	// Reading log routes
	log.Println("Registering reading log routes...")
	api.HandleFunc("/reading-logs", readingLogHandler.CreateReadingLog).Methods("POST")
	api.HandleFunc("/reading-logs", readingLogHandler.GetReadingLogs).Methods("GET")
	api.HandleFunc("/children/reading-logs", readingLogHandler.GetChildReadingLogs).Methods("GET")
	log.Println("Reading log routes registered successfully")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)

	// Add CORS middleware for development
	h := gorillahandlers.CORS(
		gorillahandlers.AllowedOrigins([]string{"*"}),
		gorillahandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		gorillahandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
	)(r)

	log.Fatal(http.ListenAndServe(":"+port, h))
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Add claims to request context
			ctx := context.WithValue(r.Context(), "user_id", uint(claims["user_id"].(float64)))
			ctx = context.WithValue(ctx, "role", claims["role"].(string))
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
	})
}

