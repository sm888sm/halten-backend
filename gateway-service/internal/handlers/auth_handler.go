package handlers

import (
	"net/http"

	pb_auth "github.com/sm888sm/halten-backend/user-service/api/pb"

	"github.com/gin-gonic/gin"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	"github.com/sm888sm/halten-backend/common/responsehandler"
	"github.com/sm888sm/halten-backend/gateway-service/internal/services"
)

type AuthHandler struct {
	services *services.Services
}

func NewAuthHandler(services *services.Services) *AuthHandler {
	return &AuthHandler{services: services}
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewHttpBadRequestError())
		return
	}

	authClient, err := h.services.GetAuthClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	loginRequest := pb_auth.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	response, err := authClient.Login(ctx, &loginRequest)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	responsehandler.Success(c, http.StatusCreated, "User logged in successfully", response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewHttpBadRequestError())
		return
	}

	authClient, err := h.services.GetAuthClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	refreshTokenRequest := pb_auth.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	response, err := authClient.RefreshToken(ctx, &refreshTokenRequest)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
