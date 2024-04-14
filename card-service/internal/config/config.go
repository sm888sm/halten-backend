package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port     int
	Database DatabaseConfig
	RabbitMQ RabbitMQConfig
	Services ServiceConfig
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

type ServiceConfig struct {
	UserServiceAddr  string
	BoardServiceAddr string
	ListServiceAddr  string
	CardServiceAddr  string
}

type RabbitMQConfig struct { // Add this struct
	URL string
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
		Services: ServiceConfig{
			UserServiceAddr:  os.Getenv("USER_SERVICE_ADDR"),
			BoardServiceAddr: os.Getenv("BOARD_SERVICE_ADDR"),
			ListServiceAddr:  os.Getenv("LIST_SERVICE_ADDR"),
		},
		RabbitMQ: RabbitMQConfig{ // Add this line
			URL: os.Getenv("RABBITMQ_URL"),
		},
	}, nil
}
