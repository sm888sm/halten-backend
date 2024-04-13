package main

import (
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	pb "github.com/sm888sm/halten-backend/board-service/api/pb"

	"github.com/sm888sm/halten-backend/board-service/internal/config"
	"github.com/sm888sm/halten-backend/board-service/internal/db"
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
	}
	sqlDB, err := db.SQLConn.DB()
	if err != nil {
		log.Fatalf("Error getting underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// Initialize repositories
	boardRepo := repositories.NewBoardRepository(db.SQLConn)

	// Initialize services
	boardService := services.NewBoardService(boardRepo)

	// Create gRPC server with validation interceptor
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middlewares.ValidationInterceptor))

	// Register services
	pb.RegisterBoardServiceServer(grpcServer, boardService)

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
