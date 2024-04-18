package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	pb "github.com/sm888sm/halten-backend/board-service/api/pb"
	external_services "github.com/sm888sm/halten-backend/board-service/external/services"
	consumer "github.com/sm888sm/halten-backend/board-service/internal/messaging/rabbitmq/consumer"

	"github.com/sm888sm/halten-backend/common/messaging/rabbitmq/publishers"

	"github.com/sm888sm/halten-backend/board-service/internal/config"
	"github.com/sm888sm/halten-backend/board-service/internal/connections/db"
	"github.com/sm888sm/halten-backend/board-service/internal/connections/rabbitmq"

	"github.com/sm888sm/halten-backend/board-service/internal/middlewares"
	"github.com/sm888sm/halten-backend/board-service/internal/repositories"
	"github.com/sm888sm/halten-backend/board-service/internal/services"
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
		log.Printf("Successfully connected to database.")
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
	boardRepo := repositories.NewBoardRepository(db.SQLConn)

	// Initialize external services
	svc := external_services.GetServices(&cfg.Services)
	defer svc.Close()

	// Initialize publishers
	publishers := &publishers.Publishers{
		CardPublisher: publishers.NewCardPublisher(rabbitmq.RabbitMQChannel),
		ListPublisher: publishers.NewListPublisher(rabbitmq.RabbitMQChannel),
	}

	// Initialize services
	boardService := services.NewBoardService(boardRepo, svc, publishers)

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
	pb.RegisterBoardServiceServer(grpcServer, boardService)

	// Run RabbitMQ Consumer
	runBoardConsumer(boardService)

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

func runBoardConsumer(boardService *services.BoardService) {
	// Get the RabbitMQ channel
	ch := rabbitmq.RabbitMQChannel

	// Initialize your consumer here.
	c := consumer.NewBoardConsumer(ch, boardService)

	// Run the consumer in a separate goroutine because it's a blocking operation
	go func() {
		err := c.ConsumeBoardMessages(context.Background())
		if err != nil {
			log.Fatalf("Failed to consume messages: %v", err)
		}
	}()
}
