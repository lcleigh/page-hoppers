package server

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"github.com/gin-contrib/cors"

	"page-hoppers-backend/internal/handlers"
)

type Server struct {
	Router            *gin.Engine
	AuthHandler       *handlers.AuthHandler
	ReadingLogHandler *handlers.ReadingLogHandler
}

func NewServer(db *gorm.DB) *Server {
	authHandler := handlers.NewAuthHandler(db, []byte(os.Getenv("JWT_SECRET")))
	readingLogHandler := handlers.NewReadingLogHandler(db)

	r := gin.New() // New router without default logger
	r.Use(gin.Logger()) // logs method, path, status, latency
	r.Use(gin.Recovery())

	// Add CORS middleware
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"}, // your frontend
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

	s := &Server{
		Router:            r,
		AuthHandler:       authHandler,
		ReadingLogHandler: readingLogHandler,
	}

	s.registerRoutes()
	// ✅ Log every route that got registered
    for _, route := range s.Router.Routes() {
        log.Printf("✅ Registered route: %-6s %s", route.Method, route.Path)
    }
	return s
}

func (s *Server) registerRoutes() {
	log.Println("Registering routes...")

	// Public routes
	s.Router.POST("/api/auth/parent/login", s.logHandler("ParentLogin", s.AuthHandler.ParentLogin))
	s.Router.POST("/api/auth/parent/register", s.logHandler("ParentRegister", s.AuthHandler.ParentRegister))
	s.Router.POST("/api/auth/child/login", s.logHandler("ChildLogin", s.AuthHandler.ChildLogin))

	// Protected routes (with JWT middleware)
	protected := s.Router.Group("/api")
	protected.Use(s.authMiddleware())

	// Children
	protected.GET("/children", s.logHandler("GetChildren", s.AuthHandler.GetChildren))
	protected.POST("/children", s.logHandler("CreateChild", s.AuthHandler.CreateChild))

	// Reading logs
	protected.POST("/reading-logs", s.logHandler("CreateReadingLog", s.ReadingLogHandler.CreateReadingLog))
	protected.GET("/reading-logs", s.logHandler("GetReadingLogs", s.ReadingLogHandler.GetReadingLogs))
	protected.GET("/children/reading-logs", s.logHandler("GetChildReadingLogs", s.ReadingLogHandler.GetChildReadingLogs))

	// Reading Summary
	protected.GET("/children/:id/summary", s.logHandler(
		"GetReadingSummary",
		s.ReadingLogHandler.GetReadingSummary,

))
}

// logHandler wraps a handler to log entry for easier debugging
func (s *Server) logHandler(name string, handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Handling route %s -> %s %s", name, c.Request.Method, c.Request.URL.Path)
		handler(c)
	}
}

func (s *Server) Start() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := s.Router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Auth middleware hit for:", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Missing Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			log.Println("Invalid token:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := uint(claims["user_id"].(float64))
			role := claims["role"].(string)
			log.Printf("Authenticated user_id=%d role=%s", userID, role)
			c.Set("user_id", userID)
			c.Set("role", role)
		} else {
			log.Println("Invalid token claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Next()
	}
}