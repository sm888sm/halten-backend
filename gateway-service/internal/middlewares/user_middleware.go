package middlewares

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	external_services "github.com/sm888sm/halten-backend/gateway-service/internal/services/external"
	pb_user "github.com/sm888sm/halten-backend/user-service/api/pb"
)

func UserMiddleware(services *external_services.Services, secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, errorhandler.NewAPIError(http.StatusUnauthorized, "Authorization header not provided"))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := validateToken(token, secretKey)
		if err != nil {
			errorhandler.HandleError(c, err)
			c.Abort()
			return
		}

		userID := (*claims)["userID"].(uint)

		userService, err := services.GetUserClient()
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
			c.Abort()
			return
		}

		user, err := userService.GetUserByID(ctx, &pb_user.GetUserByIDRequest{Id: uint64(userID)})
		if err != nil {
			errorhandler.HandleError(c, err)
			c.Abort()
			return
		}

		c.Set("user", user.User)
		c.Next()
	}
}

func validateToken(tokenString string, secretKey string) (*jwt.MapClaims, error) {
	invalidTokenError := errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, invalidTokenError
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, invalidTokenError
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_, ok := claims["userID"].(uint)
		if !ok {
			return nil, invalidTokenError
		}
		return &claims, nil
	} else {
		return nil, invalidTokenError
	}
}
