package main

import (
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/sm888sm/halten-backend/models"
	pb "github.com/sm888sm/halten-backend/user-service/api/pb"

	"github.com/sm888sm/halten-backend/user-service/internal/config"
	"github.com/sm888sm/halten-backend/user-service/internal/db"
	"github.com/sm888sm/halten-backend/user-service/internal/middlewares"
	"github.com/sm888sm/halten-backend/user-service/internal/repositories"
	"github.com/sm888sm/halten-backend/user-service/internal/services"
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

	models.Migrate(db.SQLConn)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db.SQLConn)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.SecretKey)
	userService := services.NewUserService(userRepo, cfg.BcryptCost)

	// Create gRPC server with validation interceptor
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middlewares.ValidationInterceptor))

	// Register services
	pb.RegisterAuthServiceServer(grpcServer, authService)
	pb.RegisterUserServiceServer(grpcServer, userService)

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
