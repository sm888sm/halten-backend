package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sm888sm/halten-backend/common/errorhandlers"
	pb_auth "github.com/sm888sm/halten-backend/user-service/api/pb" // Assuming your gRPC definitions are here
	"github.com/sm888sm/halten-backend/user-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	userRepo repositories.UserRepository
	pb_auth.UnimplementedAuthServiceServer
	secretKey string // Secret key used to sign JWTs
}

func NewAuthService(userRepo repositories.UserRepository, secretKey string) *AuthService {
	return &AuthService{userRepo: userRepo, secretKey: secretKey}
}

func (s *AuthService) Login(ctx context.Context, req *pb_auth.LoginRequest) (*pb_auth.LoginResponse, error) {
	res, err := s.userRepo.GetUserByUsername(
		&repositories.GetUserByUsernameRequest{
			Username: req.Username,
		})

	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(res.User.Password), []byte(req.Password)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, errorhandlers.NewAPIError(http.StatusUnauthorized, "Invalid credentials").Error())
	}

	accessToken, err := s.generateToken(res.User.ID, 15*time.Minute)
	if err != nil {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	refreshToken, err := s.generateToken(res.User.ID, 24*time.Hour)
	if err != nil {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	return &pb_auth.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *pb_auth.RefreshTokenRequest) (*pb_auth.RefreshTokenResponse, error) {
	// Validate the refresh token...
	claims, err := s.validateToken(req.RefreshToken)
	if err != nil {
		fmt.Println("Error", err)
		return nil, err
	}

	// Generate a new access token...
	accessToken, err := s.generateToken((*claims)["userID"].(uint64), 15*time.Minute)
	if err != nil {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	return &pb_auth.RefreshTokenResponse{AccessToken: accessToken}, nil
}

func (s *AuthService) CheckBoardUserRole(ctx context.Context, req *pb_auth.CheckBoardUserRoleRequest) (*pb_auth.CheckBoardUserRoleResponse, error) {
	checkReq := &repositories.CheckBoardUserRoleRequest{
		UserID:       req.UserID,
		BoardID:      req.BoardID,
		RequiredRole: req.RequiredRole,
	}

	if err := s.userRepo.CheckBoardUserRole(checkReq); err != nil {
		return nil, err
	}

	return &pb_auth.CheckBoardUserRoleResponse{Message: ""}, nil
}

func (s *AuthService) CheckBoardVisibility(ctx context.Context, req *pb_auth.CheckBoardVisibilityRequest) (*pb_auth.CheckBoardVisibilityResponse, error) {
	checkReq := &repositories.CheckBoardVisibilityRequest{
		UserID:  req.UserID,
		BoardID: req.BoardID,
	}

	if err := s.userRepo.CheckBoardVisibility(checkReq); err != nil {
		return nil, err
	}

	return &pb_auth.CheckBoardVisibilityResponse{Message: ""}, nil
}
