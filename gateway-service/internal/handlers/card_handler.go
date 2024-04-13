package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb_card "github.com/sm888sm/halten-backend/card-service/api/pb"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	"github.com/sm888sm/halten-backend/common/responsehandler"
	"github.com/sm888sm/halten-backend/gateway-service/internal/services"
	pb_user "github.com/sm888sm/halten-backend/user-service/api/pb"
)

type CreateCardInput struct {
	ListId uint64 `json:"listId"`
	Name   string `json:"name"`
}

type UpdateCardInput struct {
	Name   string `json:"name"`
	ListId uint64 `json:"listId"`
}

type MoveCardPositionInput struct {
	NewPosition int32  `json:"newPosition"`
	ListId      uint64 `json:"listId"`
}

type CardHandler struct {
	services *services.Services
}

func NewCardHandler(services *services.Services) *CardHandler {
	return &CardHandler{services: services}
}

func (h *CardHandler) CreateCard(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).Id

	var input CreateCardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewHttpBadRequestError())
		return
	}

	cardService, err := h.services.GetCardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_card.CreateCardRequest{
		UserId: userId,
		ListId: input.ListId,
		Name:   input.Name,
	}
	_, err = cardService.CreateCard(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Card created successfully", nil)
}

func (h *CardHandler) GetCardsByList(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).Id

	listIDStr := c.Param("id")
	listID, err := strconv.ParseUint(listIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid list ID"))
		return
	}

	cardService, err := h.services.GetCardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_card.GetCardsByListRequest{ListId: listID, UserId: uint64(userId)}
	resp, err := cardService.GetCardsByList(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Cards fetched successfully", resp)
}

func (h *CardHandler) UpdateCard(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).Id

	cardIDStr := c.Param("id")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid card ID"))
		return
	}

	var input UpdateCardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewHttpBadRequestError())
		return
	}

	cardService, err := h.services.GetCardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_card.UpdateCardRequest{
		Id:     cardID,
		UserId: userId,
		Name:   input.Name,
		ListId: input.ListId,
	}

	_, err = cardService.UpdateCard(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Card updated successfully", nil)
}

func (h *CardHandler) DeleteCard(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).Id

	cardIDStr := c.Param("id")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid card ID"))
		return
	}

	cardService, err := h.services.GetCardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_card.DeleteCardRequest{
		Id:     cardID,
		UserId: userId,
	}

	_, err = cardService.DeleteCard(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Card deleted successfully", nil)
}

func (h *CardHandler) MoveCardPosition(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).Id

	cardIDStr := c.Param("id")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid card ID"))
		return
	}

	var input MoveCardPositionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewHttpBadRequestError())
		return
	}

	cardService, err := h.services.GetCardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_card.MoveCardPositionRequest{
		Id:          cardID,
		UserId:      userId,
		NewPosition: input.NewPosition,
		ListId:      input.ListId,
	}

	_, err = cardService.MoveCardPosition(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Card position moved successfully", nil)
}
