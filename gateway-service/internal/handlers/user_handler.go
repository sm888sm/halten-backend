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

type CreateUserBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var body CreateUserBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, errorhandlers.NewHttpBadRequestError())
		return
	}

	userService, err := h.services.GetUserClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	grpcUserReq := &pb_user.CreateUserRequest{
		Username: body.Username,
		Password: body.Password,
		Email:    body.Email,
		Fullname: body.Fullname,
	}

	resp, err := userService.CreateUser(c, grpcUserReq)
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

type UpdateEmailBody struct {
	UserID   uint64 `json:"userID" binding:"required"`
	NewEmail string `json:"newEmail" binding:"required"`
}

func (h *UserHandler) UpdateEmail(c *gin.Context) {
	var body UpdateEmailBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, errorhandlers.NewHttpBadRequestError())
		return
	}

	userService, err := h.services.GetUserClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	grpcUserReq := &pb_user.UpdateEmailRequest{
		UserID:   body.UserID,
		NewEmail: body.NewEmail,
	}

	resp, err := userService.UpdateEmail(c, grpcUserReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Email updated successfully", resp)
}

type UpdatePasswordBody struct {
	UserID      uint64 `json:"userID" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var body UpdatePasswordBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, errorhandlers.NewHttpBadRequestError())
		return
	}

	userService, err := h.services.GetUserClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	grpcUserReq := &pb_user.UpdatePasswordRequest{
		UserID:      body.UserID,
		NewPassword: body.NewPassword,
	}

	resp, err := userService.UpdatePassword(c, grpcUserReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Password updated successfully", resp)
}

type UpdateUsernameBody struct {
	UserID   uint64 `json:"userID" binding:"required"`
	Username string `json:"username" binding:"required"`
}

func (h *UserHandler) UpdateUsername(c *gin.Context) {
	var body UpdateUsernameBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, errorhandlers.NewHttpBadRequestError())
		return
	}

	userService, err := h.services.GetUserClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	grpcUserReq := &pb_user.UpdateUsernameRequest{
		UserID:   body.UserID,
		Username: body.Username,
	}

	resp, err := userService.UpdateUsername(c, grpcUserReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Username updated successfully", resp)
}

type ConfirmEmailBody struct {
	UserID uint64 `json:"userID" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

func (h *UserHandler) ConfirmEmail(c *gin.Context) {
	var body ConfirmEmailBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, errorhandlers.NewHttpBadRequestError())
		return
	}

	userService, err := h.services.GetUserClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	grpcUserReq := &pb_user.ConfirmEmailRequest{
		UserID: body.UserID,
		Token:  body.Token,
	}

	resp, err := userService.ConfirmEmail(c, grpcUserReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandlers.NewHttpInternalError())
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Email confirmed successfully", resp)
}
