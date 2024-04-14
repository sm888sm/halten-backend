package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/sm888sm/halten-backend/gateway-service/internal/config"
	"github.com/sm888sm/halten-backend/gateway-service/internal/connections/rabbitmq"

	"github.com/sm888sm/halten-backend/gateway-service/internal/routes"
	external_services "github.com/sm888sm/halten-backend/gateway-service/internal/services/external"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Connect to RabbitMQ
	err = rabbitmq.Connect(&cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}

	// Get services
	svc := external_services.GetServices(&cfg.Services)

	defer svc.Close()

	// Initialize Gin
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, svc, cfg.SecretKey)

	// Start the Gin server
	r.Run(":" + cfg.Port)
	log.Println("Application started.")
}
