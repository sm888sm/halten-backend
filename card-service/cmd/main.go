package main

import (
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	pb_card "github.com/sm888sm/halten-backend/card-service/api/pb"
	external_services "github.com/sm888sm/halten-backend/card-service/external/services"
	"github.com/sm888sm/halten-backend/card-service/internal/middlewares"
	"github.com/sm888sm/halten-backend/card-service/internal/services"

	"github.com/sm888sm/halten-backend/card-service/internal/config"
	"github.com/sm888sm/halten-backend/card-service/internal/connections/db"
	"github.com/sm888sm/halten-backend/card-service/internal/connections/rabbitmq"
	"github.com/sm888sm/halten-backend/card-service/internal/repositories"
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

	// Connect to database
	err = db.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	} else {
		log.Printf("Successfully connected to the database.")
	}

	sqlDB, err := db.SQLConn.DB()
	if err != nil {
		log.Fatalf("Error getting underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// Connect to RabbitMQ
	err = rabbitmq.Connect(&cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	} else {
		log.Printf("Successfully connected to RabbitMQ.")
	}

	// Initialize repositories
	cardRepo := repositories.NewCardRepository(db.SQLConn)

	// Initialize external services
	svc := external_services.GetServices(&cfg.Services)
	defer svc.Close()

	// Initialize services
	cardService := services.NewCardService(cardRepo)

	// Create gRPC server with validation interceptor
	AuthInterceptor := middlewares.NewAuthInterceptor(db.SQLConn, svc)
	validatorInterceptor := middlewares.NewValidatorInterceptor(db.SQLConn)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			AuthInterceptor.AuthInterceptor,
			validatorInterceptor.ValidationInterceptor,
		),
	)

	// Register services
	pb_card.RegisterCardServiceServer(grpcServer, cardService)

	// Start listening
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Service listening on port %d", cfg.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
