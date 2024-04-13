package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/sm888sm/halten-backend/gateway-service/internal/config"
	"github.com/sm888sm/halten-backend/gateway-service/internal/routes"
	"github.com/sm888sm/halten-backend/gateway-service/internal/services"
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

	// Get services
	svc := services.GetServices(&cfg.Services)

	defer svc.Close()

	// Initialize Gin
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, svc, cfg.SecretKey)

	// Start the Gin server
	r.Run(":" + cfg.Port)
	log.Println("Application started.")
}
