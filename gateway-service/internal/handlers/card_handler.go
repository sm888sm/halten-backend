package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	pb_board "github.com/sm888sm/halten-backend/board-service/api/pb"
	pb_card "github.com/sm888sm/halten-backend/card-service/api/pb"
	"github.com/sm888sm/halten-backend/common/errorhandlers"
	"github.com/sm888sm/halten-backend/common/responsehandlers"
	external_services "github.com/sm888sm/halten-backend/gateway-service/external/services"
	pb_auth "github.com/sm888sm/halten-backend/user-service/api/pb"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CardHandler struct {
	services *external_services.Services
}

func NewCardHandler(services *external_services.Services) *CardHandler {
	return &CardHandler{services: services}
}

/*
********************
* No Authorization *
********************
 */

type GetCardByIDUri struct {
	CardID uint64 `uri:"cardID" binding:"required"`
}

func (h *CardHandler) GetCardByID(c *gin.Context) {
	ctx := c.Request.Context()

	var uri GetCardByIDUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewHttpInternalError())
		return
	}

	if err := h.CheckVisibility(ctx, userID, uri.CardID); err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewHttpInternalError())
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcReq := &pb_card.GetCardByIDRequest{CardID: uri.CardID}

	resp, err := cardClient.GetCardByID(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Card retrieved successfully", resp)
}

type GetCardsByBoardUri struct {
	BoardID uint64 `uri:"boardID" binding:"required"`
}

func (h *CardHandler) GetCardsByBoard(c *gin.Context) {
	ctx := c.Request.Context()

	var uri GetCardsByBoardUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	if err := h.CheckVisibility(ctx, userID, uri.BoardID); err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	md := metadata.Pairs("boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.GetCardsByBoardRequest{}

	res, err := cardClient.GetCardsByBoard(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Card list retrieved successfully", res.Cards)
}

type GetCardsByListUri struct {
	ListID uint64 `uri:"listID" binding:"required"`
}

func (h *CardHandler) GetCardsByList(c *gin.Context) {
	ctx := c.Request.Context()

	var uri GetCardsByListUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewHttpInternalError())
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.ListID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
	}

	boardID := grpcBoardRes.BoardID

	if err := h.CheckVisibility(ctx, userID, boardID); err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcCardReq := &pb_card.GetCardsByListRequest{ListID: uri.ListID}

	res, err := cardClient.GetCardsByList(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Card list retrieved successfully", res.Cards)
}

/*
****************************
* Authorization Required *
****************************
 */

type CreateCardRequest struct {
	Name   string `json:"name" binding:"required"`
	ListID uint64 `json:"listID" binding:"required"`
}

func (h *CardHandler) CreateCard(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByListRequest{
		ListID: req.ListID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByList(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.CreateCardRequest{
		ListID: req.ListID,
		Name:   req.Name,
	}
	grpcCardRes, err := cardClient.CreateCard(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Card created successfully", grpcCardRes.Card)
}

type MoveCardPositionUri struct {
	CardID uint64 `uri:"listID" binding:"required"`
}

type MoveCardPositionBody struct {
	Position  int64  `json:"position" binding:"required"`
	OldListID uint64 `json:"oldListID" binding:"required"`
	NewListID uint64 `json:"newListID" binding:"required"`
}

func (h *CardHandler) MoveCardPosition(c *gin.Context) {
	ctx := c.Request.Context()

	var uri MoveCardPositionUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body MoveCardPositionBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.MoveCardPositionRequest{
		CardID:    uri.CardID,
		Position:  body.Position,
		OldListID: body.OldListID,
		NewListID: body.NewListID,
	}

	grpcCardRes, err := cardClient.MoveCardPosition(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type UpdateCardNameUri struct {
	CardID uint64 `uri:"listID" binding:"required"`
}

type UpdateCardNameBody struct {
	Name string `json:"name" binding:"required"`
}

func (h *CardHandler) UpdateCardName(c *gin.Context) {
	ctx := c.Request.Context()

	var uri UpdateCardNameUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body UpdateCardNameBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.UpdateCardNameRequest{
		CardID: uri.CardID,
		Name:   body.Name,
	}

	grpcCardRes, err := cardClient.UpdateCardName(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type UpdateCardDescriptionUri struct {
	CardID uint64 `uri:"listID" binding:"required"`
}

type UpdateCardDescriptionBody struct {
	Description string `json:"description" binding:"required"`
}

func (h *CardHandler) UpdateCardDescription(c *gin.Context) {
	ctx := c.Request.Context()

	var uri UpdateCardDescriptionUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body UpdateCardDescriptionBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.UpdateCardDescriptionRequest{
		CardID:      uri.CardID,
		Description: body.Description,
	}

	grpcCardRes, err := cardClient.UpdateCardDescription(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type AddCardLabelUri struct {
	CardID  uint64 `uri:"listID" binding:"required"`
	LabelID uint64 `uri:"labelID" binding:"required"`
}

func (h *CardHandler) AddCardLabel(c *gin.Context) {
	ctx := c.Request.Context()

	var uri AddCardLabelUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.AddCardLabelRequest{
		CardID:  uri.CardID,
		LabelID: uri.LabelID,
	}

	grpcCardRes, err := cardClient.AddCardLabel(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type RemoveCardLabelUri struct {
	CardID  uint64 `uri:"listID" binding:"required"`
	LabelID uint64 `uri:"labelID" binding:"required"`
}

func (h *CardHandler) RemoveCardLabel(c *gin.Context) {
	ctx := c.Request.Context()

	var uri RemoveCardLabelUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.RemoveCardLabelRequest{
		CardID:  uri.CardID,
		LabelID: uri.LabelID,
	}

	grpcCardRes, err := cardClient.RemoveCardLabel(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type SetCardDatesUri struct {
	CardID uint64 `uri:"listID" binding:"required"`
}

type SetCardDatesBody struct {
	StartDate string `json:"startDate"`
	DueDate   string `json:"dueDate"`
}

func (h *CardHandler) SetCardDates(c *gin.Context) {
	ctx := c.Request.Context()

	var uri SetCardDatesUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body SetCardDatesBody
	if err := c.ShouldBind(&body); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	var startDate *timestamppb.Timestamp
	var dueDate *timestamppb.Timestamp
	var fieldErrors []errorhandlers.FieldError

	if body.StartDate != "" {
		parsedStartDate, err := time.Parse(time.RFC3339, body.StartDate)
		if err != nil {
			fieldErrors = append(fieldErrors, errorhandlers.FieldError{
				Code:    "InvalidDateFormat",
				Message: "Invalid start date format",
				Field:   "startDate",
			})
		} else {
			startDate = timestamppb.New(parsedStartDate)
		}
	}

	if body.DueDate != "" {
		parsedDueDate, err := time.Parse(time.RFC3339, body.DueDate)
		if err != nil {
			fieldErrors = append(fieldErrors, errorhandlers.FieldError{
				Code:    "InvalidDateFormat",
				Message: "Invalid due date format",
				Field:   "dueDate",
			})
		} else {
			dueDate = timestamppb.New(parsedDueDate)
		}
	}

	if len(fieldErrors) > 0 {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid request parameters", fieldErrors...))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.SetCardDatesRequest{
		CardID:    uri.CardID,
		StartDate: startDate,
		DueDate:   dueDate,
	}

	grpcCardRes, err := cardClient.SetCardDates(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type ToggleCardCompletedUri struct {
	CardID uint64 `uri:"listID" binding:"required"`
}

func (h *CardHandler) ToggleCardCompleted(c *gin.Context) {
	ctx := c.Request.Context()

	var uri ToggleCardCompletedUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.ToggleCardCompletedRequest{
		CardID: uri.CardID,
	}

	grpcCardRes, err := cardClient.ToggleCardCompleted(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type AddCardAttachmentUri struct {
	CardID       uint64 `uri:"listID" binding:"required"`
	AttachmentID uint64 `uri:"attachmentID" binding:"required"`
}

func (h *CardHandler) AddCardAttachment(c *gin.Context) {
	ctx := c.Request.Context()

	var uri AddCardAttachmentUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.AddCardAttachmentRequest{
		CardID:       uri.CardID,
		AttachmentID: uri.AttachmentID,
	}

	grpcCardRes, err := cardClient.AddCardAttachment(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type RemoveCardAttachmentUri struct {
	CardID       uint64 `uri:"listID" binding:"required"`
	AttachmentID uint64 `uri:"attachmentID" binding:"required"`
}

func (h *CardHandler) RemoveCardAttachment(c *gin.Context) {
	ctx := c.Request.Context()

	var uri RemoveCardAttachmentUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.RemoveCardAttachmentRequest{
		CardID:       uri.CardID,
		AttachmentID: uri.AttachmentID,
	}

	grpcCardRes, err := cardClient.RemoveCardAttachment(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type AddCardCommentUri struct {
	CardID uint64 `uri:"listID" binding:"required"`
}

type AddCardCommentBody struct {
	Content string `uri:"content" binding:"required"`
}

func (h *CardHandler) AddCardComment(c *gin.Context) {
	ctx := c.Request.Context()

	var uri AddCardCommentUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body AddCardCommentBody
	if err := c.ShouldBind(&body); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.AddCardCommentRequest{
		CardID:  uri.CardID,
		Content: body.Content,
	}

	grpcCardRes, err := cardClient.AddCardComment(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type RemoveCardCommentUri struct {
	CardID    uint64 `uri:"listID" binding:"required"`
	CommendID uint64 `uri:"commentID" binding:"required"`
}

func (h *CardHandler) RemoveCardComment(c *gin.Context) {
	ctx := c.Request.Context()

	var uri RemoveCardCommentUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.RemoveCardCommentRequest{
		CardID:    uri.CardID,
		CommentID: uri.CommendID,
	}

	grpcCardRes, err := cardClient.RemoveCardComment(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type AddCardMembersUri struct {
	CardID uint64 `uri:"cardID" binding:"required"`
}

type AddCardMembersBody struct {
	UserIDs []uint64 `json:"userIDs" binding:"required"`
}

func (h *CardHandler) AddCardMembers(c *gin.Context) {
	ctx := c.Request.Context()

	var uri AddCardMembersUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body AddCardMembersBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.AddCardMembersRequest{
		CardID:  uri.CardID,
		UserIDs: body.UserIDs,
	}

	grpcCardRes, err := cardClient.AddCardMembers(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type RemoveCardMembersUri struct {
	CardID uint64 `uri:"cardID" binding:"required"`
}

type RemoveCardMembersBody struct {
	UserIDs []uint64 `json:"userIDs" binding:"required"`
}

func (h *CardHandler) RemoveCardMembers(c *gin.Context) {
	ctx := c.Request.Context()

	var uri RemoveCardMembersUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body RemoveCardMembersBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.RemoveCardMembersRequest{
		CardID:  uri.CardID,
		UserIDs: body.UserIDs,
	}

	grpcCardRes, err := cardClient.RemoveCardMembers(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type ArchiveCardUri struct {
	CardID uint64 `uri:"cardID" binding:"required"`
}

func (h *CardHandler) ArchiveCard(c *gin.Context) {
	ctx := c.Request.Context()

	var uri ArchiveCardUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardResp, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardResp.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.ArchiveCardRequest{
		CardID: uri.CardID,
	}

	grpcCardRes, err := cardClient.ArchiveCard(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type RestoreCardUri struct {
	CardID uint64 `uri:"cardID" binding:"required"`
}

func (h *CardHandler) RestoreCard(c *gin.Context) {
	ctx := c.Request.Context()

	var uri RestoreCardUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardResp, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardResp.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.RestoreCardRequest{CardID: uri.CardID}

	grpcCardRes, err := cardClient.RestoreCard(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

type DeleteCardUri struct {
	CardID uint64 `uri:"cardID" binding:"required"`
}

func (h *CardHandler) DeleteCard(c *gin.Context) {
	ctx := c.Request.Context()

	var uri DeleteCardUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	cardClient, err := h.services.GetCardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByCardRequest{
		CardID: uri.CardID,
	}

	grpcBoardResp, err := boardClient.GetBoardIDByCard(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardResp.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcCardReq := &pb_card.DeleteCardRequest{
		CardID: uri.CardID,
	}

	grpcCardRes, err := cardClient.DeleteCard(ctx, grpcCardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcCardRes.Message, nil)
}

// Helpers

func (h *CardHandler) CheckVisibility(ctx context.Context, userID, boardID uint64) error {
	authClient, err := h.services.GetAuthClient()
	if err != nil {
		return err
	}

	_, err = authClient.CheckBoardVisibility(ctx, &pb_auth.CheckBoardVisibilityRequest{
		UserID:  userID,
		BoardID: boardID,
	})

	return err
}
