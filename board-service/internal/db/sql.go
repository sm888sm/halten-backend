package db

import (
	"fmt"

	"github.com/sm888sm/halten-backend/board-service/internal/config"

	"gorm.io/driver/postgres" // Import the PostgreSQL driver
	"gorm.io/gorm"
)

var SQLConn *gorm.DB

func Connect(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.Username, cfg.Password, cfg.DBName, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	SQLConn = db
	return nil
}
