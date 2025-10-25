package server

import (
    "context"
    "log"
    "net/http"
    "os"
    "strings"

    "github.com/golang-jwt/jwt/v5"
    "github.com/gorilla/mux"
    gorillahandlers "github.com/gorilla/handlers"
    "page-hoppers-backend/internal/handlers"
    "gorm.io/gorm"
)

type Server struct {
    Router          *mux.Router
    AuthHandler     *handlers.AuthHandler
    ReadingLogHandler *handlers.ReadingLogHandler
}

func NewServer(db *gorm.DB) *Server {
    authHandler := handlers.NewAuthHandler(db, []byte(os.Getenv("JWT_SECRET")))
    readingLogHandler := handlers.NewReadingLogHandler(db)

    s := &Server{
        Router: mux.NewRouter(),
        AuthHandler: authHandler,
        ReadingLogHandler: readingLogHandler,
    }

    s.registerRoutes()
    return s
}

func (s *Server) registerRoutes() {
    // Public routes
    s.Router.HandleFunc("/api/auth/parent/login", s.AuthHandler.ParentLogin).Methods("POST")
    s.Router.HandleFunc("/api/auth/parent/register", s.AuthHandler.ParentRegister).Methods("POST")
    s.Router.HandleFunc("/api/auth/child/login", s.AuthHandler.ChildLogin).Methods("POST")

    // Protected routes
    api := s.Router.PathPrefix("/api").Subrouter()
    api.Use(authMiddleware)

    api.HandleFunc("/children", s.AuthHandler.GetChildren).Methods("GET")
    api.HandleFunc("/children", s.AuthHandler.CreateChild).Methods("POST")

    api.HandleFunc("/reading-logs", s.ReadingLogHandler.CreateReadingLog).Methods("POST")
    api.HandleFunc("/reading-logs", s.ReadingLogHandler.GetReadingLogs).Methods("GET")
    api.HandleFunc("/children/reading-logs", s.ReadingLogHandler.GetChildReadingLogs).Methods("GET")
}

func (s *Server) Start() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)

    handler := gorillahandlers.CORS(
        gorillahandlers.AllowedOrigins([]string{"*"}),
        gorillahandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
        gorillahandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
    )(s.Router)

    log.Fatal(http.ListenAndServe(":"+port, handler))
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