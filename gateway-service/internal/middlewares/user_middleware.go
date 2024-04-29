package middlewares

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/errorhandlers"
	external_services "github.com/sm888sm/halten-backend/gateway-service/external/services"
	pb_user "github.com/sm888sm/halten-backend/user-service/api/pb"
)

func UserMiddleware(services *external_services.Services, secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, errorhandlers.NewAPIError(http.StatusUnauthorized, "Authorization header not provided"))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := validateToken(token, secretKey)
		if err != nil {
			errorhandlers.HandleError(c, err)
			c.Abort()
			return
		}

		userID := (*claims)["userID"].(uint64)

		userService, err := services.GetUserClient()
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
			c.Abort()
			return
		}

		_, err = userService.GetUserByID(ctx, &pb_user.GetUserByIDRequest{UserID: uint64(userID)})
		if err != nil {
			errorhandlers.HandleError(c, err)
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func validateToken(tokenString string, secretKey string) (*jwt.MapClaims, error) {
	invalidTokenError := errorhandlers.NewAPIError(httpcodes.ErrBadRequest, "Invalid token")

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
