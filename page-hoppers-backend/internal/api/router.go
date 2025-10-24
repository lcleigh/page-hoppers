package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) http.Handler {
	r := mux.NewRouter()

	// wire up routes e.g.
	// authHandler := handlers.NewAuthHandler(db, []byte(os.Getenv("JWT_SECRET")))
	// r.HandleFunc("/api/auth/parent/login", authHandler.ParentLogin).Methods("POST")

	return r
}