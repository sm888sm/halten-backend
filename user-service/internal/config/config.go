package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port       int
	Database   DatabaseConfig
	SecretKey  string
	BcryptCost int
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
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8000 // Default port
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		dbPort = 5432 // Default PostgreSQL port
	}

	bcryptCost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		bcryptCost = 10 // Default bcrypt cost
	}

	return &Config{
		Port: port, // Or your default
		Database: DatabaseConfig{
			Driver:   "postgres", // Changed to postgres
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			DBName:   os.Getenv("DB_NAME"),
		},
		SecretKey:  os.Getenv("SECRET_KEY"),
		BcryptCost: bcryptCost,
	}, nil
}
