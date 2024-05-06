package handlers

import (
	"context"
	"net/http"
	"strconv"

	pb_board "github.com/sm888sm/halten-backend/board-service/api/pb"
	pb_list "github.com/sm888sm/halten-backend/list-service/api/pb"
	pb_auth "github.com/sm888sm/halten-backend/user-service/api/pb"
	"google.golang.org/grpc/metadata"

	"github.com/gin-gonic/gin"
	"github.com/sm888sm/halten-backend/common/errorhandlers"
	"github.com/sm888sm/halten-backend/common/responsehandlers"
	external_services "github.com/sm888sm/halten-backend/gateway-service/external/services"
)

type ListHandler struct {
	services *external_services.Services
}

func NewListHandler(services *external_services.Services) *ListHandler {
	return &ListHandler{services: services}
}

/*
********************
* No Authorization *
********************
 */

type GetListByIDUri struct {
	ListID uint64 `uri:"listID" binding:"required"`
}

func (h *ListHandler) GetListByID(c *gin.Context) {
	ctx := c.Request.Context()

	var uri GetListByIDUri
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

	listClient, err := h.services.GetListClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByListRequest{
		ListID: uri.ListID,
	}

	grpcBoardResp, err := boardClient.GetBoardIDByList(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardResp.BoardID

	if err := h.CheckVisibility(c.Request.Context(), userID, boardID); err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	md := metadata.Pairs("boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcListReq := &pb_list.GetListByIDRequest{ListID: uri.ListID} // Convert the HTTP request to the gRPC request

	grpcListResp, err := listClient.GetListByID(ctx, grpcListReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, "List retrieved successfully", grpcListResp.Lists)
}

type GetListsByBoardUri struct {
	BoardID uint64 `uri:"boardID" binding:"required"`
}

func (h *ListHandler) GetListsByBoard(c *gin.Context) {
	ctx := c.Request.Context()

	var uri GetListsByBoardUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	listClient, err := h.services.GetListClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	if err := h.CheckVisibility(ctx, userID, uri.BoardID); err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	md := metadata.Pairs("boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcReq := &pb_list.GetListsByBoardRequest{} // Convert the HTTP request to the gRPC request

	resp, err := listClient.GetListsByBoard(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, "List retrieved successfully", resp)
}

/*
****************************
* Authorization Required *
****************************
 */

type CreateListBody struct {
	Name    string `json:"name" binding:"required"`
	BoardID uint64 `json:"boardID" binding:"required"`
}

func (h *ListHandler) CreateList(c *gin.Context) {
	var body CreateListBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	listClient, err := h.services.GetListClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(body.BoardID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcListReq := &pb_list.CreateListRequest{
		Name: body.Name,
	}
	grpcListRes, err := listClient.CreateList(ctx, grpcListReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, "List created successfully", grpcListRes.List)
}

type MoveListPositionUri struct {
	ListID uint64 `uri:"listID" binding:"required"`
}

type MoveListPositionBody struct {
	Position int64 `json:"position" binding:"required"`
}

func (h *ListHandler) MoveListPosition(c *gin.Context) {
	ctx := c.Request.Context()

	var uri MoveListPositionUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body MoveListPositionBody
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

	listClient, err := h.services.GetListClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByListRequest{
		ListID: uri.ListID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByList(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcListReq := &pb_list.MoveListPositionRequest{ListID: uri.ListID, Position: body.Position}

	grpcListRes, err := listClient.MoveListPosition(ctx, grpcListReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcListRes.Message, nil)
}

type UpdateListNameUri struct {
	ListID uint64 `uri:"listID" binding:"required"`
}

type UpdateListNameBody struct {
	Name string `json:"name" binding:"required"`
}

func (h *ListHandler) UpdateListName(c *gin.Context) {
	ctx := c.Request.Context()

	var uri UpdateListNameUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body UpdateListNameBody
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

	listClient, err := h.services.GetListClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByListRequest{
		ListID: uri.ListID,
	}

	grpcBoardRes, err := boardClient.GetBoardIDByList(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardRes.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcListReq := &pb_list.UpdateListNameRequest{
		ListID: uri.ListID,
		Name:   body.Name,
	}

	grpcListRes, err := listClient.UpdateListName(ctx, grpcListReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcListRes.Message, nil)
}

type ArchiveListUri struct {
	ListID uint64 `uri:"listID" binding:"required"`
}

func (h *ListHandler) ArchiveList(c *gin.Context) {
	ctx := c.Request.Context()

	var uri ArchiveListUri
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

	listClient, err := h.services.GetListClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByListRequest{
		ListID: uri.ListID,
	}

	grpcBoardResp, err := boardClient.GetBoardIDByList(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardResp.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcListReq := &pb_list.ArchiveListRequest{ListID: uri.ListID}

	grpcListRes, err := listClient.ArchiveList(ctx, grpcListReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcListRes.Message, nil)
}

type RestoreListUri struct {
	ListID uint64 `uri:"listID" binding:"required"`
}

func (h *ListHandler) RestoreList(c *gin.Context) {
	ctx := c.Request.Context()

	var uri RestoreListUri
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

	listClient, err := h.services.GetListClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByListRequest{
		ListID: uri.ListID,
	}

	grpcBoardResp, err := boardClient.GetBoardIDByList(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardResp.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcListReq := &pb_list.RestoreListRequest{ListID: uri.ListID}

	grpcListRes, err := listClient.RestoreList(ctx, grpcListReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcListRes.Message, nil)
}

type DeleteListUri struct {
	ListID uint64 `uri:"listID" binding:"required"`
}

func (h *ListHandler) DeleteList(c *gin.Context) {
	ctx := c.Request.Context()

	var uri DeleteListUri
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

	listClient, err := h.services.GetListClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	grpcBoardReq := &pb_board.GetBoardIDByListRequest{
		ListID: uri.ListID,
	}

	grpcBoardResp, err := boardClient.GetBoardIDByList(ctx, grpcBoardReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardID := grpcBoardResp.BoardID

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(boardID, 10))
	ctx = metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcListReq := &pb_list.DeleteListRequest{}

	grpcListRes, err := listClient.DeleteList(ctx, grpcListReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, grpcListRes.Message, nil)
}

// Helpers

func (h *ListHandler) CheckVisibility(ctx context.Context, userID, boardID uint64) error {
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
