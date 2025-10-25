package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"page-hoppers-backend/internal/repository"
	"page-hoppers-backend/internal/server"
)

func main() {
	if os.Getenv("IN_DOCKER") != "true" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found (this is fine if running in Docker)")
		}
	}

	db := repository.InitDB()
	srv := server.NewServer(db)
	srv.Start()
}