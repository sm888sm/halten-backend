package services

import (
	"time"

	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/errorhandler" // Assuming your gRPC definitions are here
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dgrijalva/jwt-go"
)

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
