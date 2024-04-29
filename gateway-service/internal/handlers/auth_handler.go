package handlers

import (
	"net/http"

	external_services "github.com/sm888sm/halten-backend/gateway-service/external/services"
	pb_auth "github.com/sm888sm/halten-backend/user-service/api/pb"

	"github.com/gin-gonic/gin"
	"github.com/sm888sm/halten-backend/common/errorhandlers"
	"github.com/sm888sm/halten-backend/common/responsehandlers"
)

type AuthHandler struct {
	services *external_services.Services
}

func NewAuthHandler(services *external_services.Services) *AuthHandler {
	return &AuthHandler{services: services}
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorhandlers.NewHttpBadRequestError())
		return
	}

	authClient, err := h.services.GetAuthClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	loginRequest := pb_auth.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	response, err := authClient.Login(ctx, &loginRequest)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusCreated, "User logged in successfully", response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorhandlers.NewHttpBadRequestError())
		return
	}

	authClient, err := h.services.GetAuthClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	refreshTokenRequest := pb_auth.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	response, err := authClient.RefreshToken(ctx, &refreshTokenRequest)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
