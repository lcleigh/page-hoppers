package integration_respository

import (
    "log"

    "page-hoppers-backend/internal/repository"
    "gorm.io/gorm"
)

func SetupTestDB() *gorm.DB {
    db := repository.InitDB()
    if db == nil {
        log.Fatal("failed to connect to test database")
    }

    return db
}