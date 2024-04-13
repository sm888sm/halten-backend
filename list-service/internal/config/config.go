package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port     int
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

func LoadConfig() (*Config, error) {
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		dbPort = 5432 // Default PostgreSQL port
	}

	return &Config{
		Port: 8080, // Or your default
		Database: DatabaseConfig{
			Driver:   "postgres", // Changed to postgres
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			DBName:   os.Getenv("DB_NAME"),
		},
	}, nil
}
