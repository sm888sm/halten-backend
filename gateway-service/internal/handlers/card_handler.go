package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb_card "github.com/sm888sm/halten-backend/card-service/api/pb"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	"github.com/sm888sm/halten-backend/common/responsehandler"
	external_services "github.com/sm888sm/halten-backend/gateway-service/external/services"
	"google.golang.org/grpc/metadata"
)

type CreateCardInput struct {
	Name   string `json:"name"`
	ListID uint64 `json:"listID"`
}

type MoveCardPositionInput struct {
	NewPosition int64  `json:"newPosition"`
	OldListID   uint64 `json:"oldlistID"`
	NewListID   uint64 `json:"newListID"`
}

type CardHandler struct {
	services *external_services.Services
}

func NewCardHandler(services *external_services.Services) *CardHandler {
	return &CardHandler{services: services}
}

func (h *CardHandler) CreateCard(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var input CreateCardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewHttpBadRequestError())
		return
	}

	// TODO : Get board id by list id

	md := metadata.Pairs("userID", fmt.Sprintf("%d", userID))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	cardService, err := h.services.GetCardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_card.CreateCardRequest{
		ListID: input.ListID,
		Name:   input.Name,
	}
	_, err = cardService.CreateCard(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Card created successfully", nil)
}

func (h *CardHandler) GetCardByID(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
	}

	// TODO get board id by card id

	md := metadata.Pairs("userID", fmt.Sprintf("%d", userID))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	cardIDStr := c.Param("card-id")
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

	req := &pb_card.GetCardByIDRequest{CardID: cardID}
	resp, err := cardService.GetCardByID(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Card fetched successfully", resp)
}

func (h *CardHandler) GetCardsByBoard(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
	}

	// TODO : get board id by card id

	md := metadata.Pairs("userID", fmt.Sprintf("%d", userID))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	boardIDStr := c.Param("board-id")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid list ID"))
		return
	}

	cardService, err := h.services.GetCardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_card.GetCardsByBoardRequest{BoardID: boardID}
	resp, err := cardService.GetCardsByBoard(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Cards fetched successfully", resp)
}

func (h *CardHandler) GetCardsByList(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
	}

	// TODO get board id by list id

	md := metadata.Pairs("userID", fmt.Sprintf("%d", userID))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	listIDStr := c.Param("list-id")
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

	req := &pb_card.GetCardsByListRequest{ListID: listID}
	resp, err := cardService.GetCardsByList(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Cards fetched successfully", resp)
}

func (h *CardHandler) MoveCardPosition(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
	}

	// TODO get board id by list id for old and new list

	md := metadata.Pairs("userID", fmt.Sprintf("%d", userID))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	cardIDStr := c.Param("card-id")
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
		CardID:      cardID,
		NewPosition: input.NewPosition,
		OldListID:   input.OldListID,
		NewListID:   input.NewListID,
	}

	_, err = cardService.MoveCardPosition(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Card position moved successfully", nil)
}

func (h *CardHandler) DeleteCard(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	// TODO get board by card id

	md := metadata.Pairs("userID", fmt.Sprintf("%d", userID))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	cardIDStr := c.Param("card-id")
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
		CardID: cardID,
	}

	_, err = cardService.DeleteCard(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Card deleted successfully", nil)
}
