package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	pb_auth "github.com/sm888sm/halten-backend/user-service/api/pb" // Assuming your gRPC definitions are here
	"github.com/sm888sm/halten-backend/user-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	userRepo repositories.UserRepository
	pb_auth.UnimplementedAuthServiceServer
	secretKey string // Secret key used to sign JWTs
}

func NewAuthService(userRepo repositories.UserRepository, secretKey string) *AuthService {
	return &AuthService{userRepo: userRepo, secretKey: secretKey}
}

func (s *AuthService) Login(ctx context.Context, in *pb_auth.LoginRequest) (*pb_auth.LoginResponse, error) {
	user, err := s.userRepo.GetUserByUsername(in.Username)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, errorhandler.NewAPIError(httpcodes.ErrUnauthorized, "Invalid credentials").Error())
	}

	accessToken, err := s.generateToken(user.ID, 15*time.Minute)
	if err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	refreshToken, err := s.generateToken(user.ID, 24*time.Hour)
	if err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	return &pb_auth.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, in *pb_auth.RefreshTokenRequest) (*pb_auth.RefreshTokenResponse, error) {
	// Validate the refresh token...
	claims, err := s.validateToken(in.RefreshToken)
	if err != nil {
		fmt.Println("Error", err)
		return nil, err
	}

	// Generate a new access token...
	accessToken, err := s.generateToken((*claims)["userID"].(uint), 15*time.Minute)
	if err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	return &pb_auth.RefreshTokenResponse{AccessToken: accessToken}, nil
}

func (s *AuthService) generateToken(userID uint, duration time.Duration) (string, error) {
	// Generate a JWT with the specified duration and the user's username as a claim...
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(duration).Unix(),
	})

	// Sign the token with the secret key...
	return token.SignedString([]byte(s.secretKey))
}

func (s *AuthService) validateToken(tokenString string) (*jwt.MapClaims, error) {
	invalidTokenError := status.Errorf(codes.InvalidArgument, errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid token").Error())

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, invalidTokenError
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, invalidTokenError
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	} else {
		return nil, invalidTokenError
	}
}
