package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sm888sm/halten-backend/common/errorhandlers"
	"github.com/sm888sm/halten-backend/common/responsehandlers"
	external_services "github.com/sm888sm/halten-backend/gateway-service/external/services"
	pb_user "github.com/sm888sm/halten-backend/user-service/api/pb"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	services *external_services.Services
}

func NewUserHandler(services *external_services.Services) *UserHandler {
	return &UserHandler{services: services}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req pb_user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorhandlers.NewHttpBadRequestError())
		return
	}

	userService, err := h.services.GetUserClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	resp, err := userService.CreateUser(c, &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
			return
		}

		var apiError errorhandlers.APIError
		if err := json.Unmarshal([]byte(st.Message()), &apiError); err != nil {
			c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
			return
		}
		c.JSON(apiError.Meta.Status, apiError)
		return
	}

	responsehandlers.Success(c, http.StatusCreated, "User created successfully", resp)
}

func (h *UserHandler) UpdateEmail(c *gin.Context) {
	var req pb_user.UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorhandlers.NewHttpBadRequestError())
		return
	}

	userService, err := h.services.GetUserClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	resp, err := userService.UpdateEmail(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Email updated successfully", resp)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var req pb_user.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorhandlers.NewHttpBadRequestError())
		return
	}

	userService, err := h.services.GetUserClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	resp, err := userService.UpdatePassword(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Password updated successfully", resp)
}

func (h *UserHandler) UpdateUsername(c *gin.Context) {
	var req pb_user.UpdateUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorhandlers.NewHttpBadRequestError())
		return
	}

	userService, err := h.services.GetUserClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	resp, err := userService.UpdateUsername(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Username updated successfully", resp)
}

func (h *UserHandler) ConfirmEmail(c *gin.Context) {
	var req pb_user.ConfirmEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorhandlers.NewHttpBadRequestError())
		return
	}

	userService, err := h.services.GetUserClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	resp, err := userService.ConfirmEmail(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Email confirmed successfully", resp)
}
