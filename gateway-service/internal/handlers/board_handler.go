package handlers

import (
	"context"
	"net/http"
	"strconv"

	pb_board "github.com/sm888sm/halten-backend/board-service/api/pb"
	pb_auth "github.com/sm888sm/halten-backend/user-service/api/pb"
	"google.golang.org/grpc/metadata"

	"github.com/gin-gonic/gin"
	"github.com/sm888sm/halten-backend/common/errorhandlers"
	"github.com/sm888sm/halten-backend/common/responsehandlers"
	external_services "github.com/sm888sm/halten-backend/gateway-service/external/services"
)

type BoardHandler struct {
	services *external_services.Services
}

func NewBoardHandler(services *external_services.Services) *BoardHandler {
	return &BoardHandler{services: services}
}

/*
********************
* No Authorization *
********************
 */

type CreateBoardRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *BoardHandler) CreateBoard(c *gin.Context) {
	var req CreateBoardRequest
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

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcReq := &pb_board.CreateBoardRequest{Name: req.Name} // Convert the HTTP request to the gRPC request
	res, err := boardClient.CreateBoard(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Board created successfully", res)
}

type GetBoardByIDUri struct {
	BoardID uint64 `uri:"boardID" binding:"required"`
}

func (h *BoardHandler) GetBoardByID(c *gin.Context) {
	var uri GetBoardByIDUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewHttpInternalError())
		return
	}

	if err := h.CheckVisibility(c.Request.Context(), userID, uri.BoardID); err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardService, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewHttpInternalError())
		return
	}

	md := metadata.Pairs("boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcReq := &pb_board.GetBoardByIDRequest{} // Convert the HTTP request to the gRPC request

	resp, err := boardService.GetBoardByID(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Board retrieved successfully", resp)
}

type GetBoardListQuery struct {
	PageNumber uint64 `form:"pageNumber,default=1"`
	PageSize   uint64 `form:"pageSize,default=10"`
}

func (h *BoardHandler) GetBoardList(c *gin.Context) {
	var query GetBoardListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid request query"))
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

	req := &pb_board.GetBoardListRequest{
		PageNumber: query.PageNumber,
		PageSize:   query.PageSize,
	}

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	res, err := boardClient.GetBoardList(ctx, req)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.SuccessWithPagination(c, http.StatusOK, "Board list retrieved successfully", res.Boards, res.Pagination)
}

type GetBoardMembersUri struct {
	BoardID uint64 `json:"boardID" binding:"required"`
}

func (h *BoardHandler) GetBoardMembers(c *gin.Context) {
	var uri GetBoardMembersUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	if err := h.CheckVisibility(c.Request.Context(), userID, uri.BoardID); err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	boardClient, err := h.services.GetBoardClient()
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	req := &pb_board.GetBoardMembersRequest{}

	md := metadata.Pairs("boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	res, err := boardClient.GetBoardMembers(ctx, req)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Board members retrieved successfully", res)
}

type GetArchivedBoardListQuery struct {
	PageNumber uint64 `form:"pageNumber,default=1"`
	PageSize   uint64 `form:"pageSize,default=10"`
}

func (h *BoardHandler) GetArchivedBoardList(c *gin.Context) {
	var query GetArchivedBoardListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid request query"))
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

	req := &pb_board.GetArchivedBoardListRequest{
		PageNumber: query.PageNumber,
		PageSize:   query.PageSize,
	}

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	res, err := boardClient.GetArchivedBoardList(ctx, req)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.SuccessWithPagination(c, http.StatusOK, "Archived board list retrieved successfully", res.Boards, res.Pagination)
}

/*
****************************
* Authorization Required *
****************************
 */

type UpdateBoardNameUri struct {
	BoardID uint64 `uri:"boardID" binding:"required"`
}

type UpdateBoardNameBody struct {
	Name string `json:"name" binding:"required"`
}

func (h *BoardHandler) UpdateBoardName(c *gin.Context) {
	var uri UpdateBoardNameUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body UpdateBoardNameBody
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

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcReq := &pb_board.UpdateBoardNameRequest{Name: body.Name} // Convert the HTTP request to the gRPC request

	res, err := boardClient.UpdateBoardName(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, res.Message, nil)
}

type AddBoardUsersUri struct {
	BoardID uint64 `uri:"boardID" binding:"required"`
}

type AddBoardUsersBody struct {
	UserIDs []uint64 `json:"user_ids" binding:"required"`
	Role    string   `json:"role" binding:"required"`
}

func (h *BoardHandler) AddBoardUsers(c *gin.Context) {
	var uri AddBoardUsersUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body AddBoardUsersBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	// Get the user ID
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

	// Create metadata with userID and boardID
	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	// Create the gRPC request
	grpcReq := &pb_board.AddBoardUsersRequest{
		UserIDs: body.UserIDs,
		Role:    body.Role,
	}

	// Use the new context with metadata
	res, err := boardClient.AddBoardUsers(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, res.Message, nil)
}

type RemoveBoardUsersUri struct {
	BoardID uint64 `uri:"boardID" binding:"required"`
}

type RemoveBoardUsersBody struct {
	UserIDs []uint64 `json:"user_ids" binding:"required"`
}

func (h *BoardHandler) RemoveBoardUsers(c *gin.Context) {
	var uri RemoveBoardUsersUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body RemoveBoardUsersBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	// Get the user ID
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

	// Create metadata with userID and boardID
	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	// Create the gRPC request
	grpcReq := &pb_board.RemoveBoardUsersRequest{
		UserIDs: body.UserIDs,
	}

	// Use the new context with metadata
	res, err := boardClient.RemoveBoardUsers(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, res.Message, nil)
}

type AssignBoardUsersRoleUri struct {
	BoardID uint64 `uri:"boardID" binding:"required"`
}

type AssignBoardUsersRoleBody struct {
	UserIDs []uint64 `json:"user_ids" binding:"required"`
	Role    string   `json:"role" binding:"required"`
}

func (h *BoardHandler) AssignBoardUsersRole(c *gin.Context) {
	var uri AssignBoardUsersRoleUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body AssignBoardUsersRoleBody
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

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcReq := &pb_board.AssignBoardUsersRoleRequest{
		UserIDs: body.UserIDs,
		Role:    body.Role,
	}

	res, err := boardClient.AssignBoardUsersRole(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, res.Message, nil)
}

type ChangeBoardOwnerUri struct {
	BoardID uint64 `uri:"boardID" binding:"required"`
}

type ChangeBoardOwnerBody struct {
	NewOwnerID uint64 `json:"new_owner_id" binding:"required"`
}

func (h *BoardHandler) ChangeBoardOwner(c *gin.Context) {
	var uri ChangeBoardOwnerUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body ChangeBoardOwnerBody
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

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcReq := &pb_board.ChangeBoardOwnerRequest{NewOwnerID: body.NewOwnerID}

	res, err := boardClient.ChangeBoardOwner(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, res.Message, nil)
}

type ChangeBoardVisibilityUri struct {
	BoardID uint64 `uri:"boardID" binding:"required"`
}

type ChangeBoardVisibilityBody struct {
	Visibility string `json:"visibility" binding:"required"`
}

func (h *BoardHandler) ChangeBoardVisibility(c *gin.Context) {
	var uri ChangeBoardVisibilityUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body ChangeBoardVisibilityBody
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

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcReq := &pb_board.ChangeBoardVisibilityRequest{Visibility: body.Visibility}

	res, err := boardClient.ChangeBoardVisibility(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, res.Message, nil)
}

type AddLabelUri struct {
	BoardID uint64 `uri:"boardID" binding:"required"`
}

type AddLabelBody struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

func (h *BoardHandler) AddLabel(c *gin.Context) {
	var uri AddLabelUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body AddLabelBody
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

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcReq := &pb_board.AddLabelRequest{
		Name:  body.Name,
		Color: body.Color,
	}

	res, err := boardClient.AddLabel(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, "Label added successfully", res.Label)
}

type RemoveLabelUri struct {
	BoardID uint64 `uri:"boardID" binding:"required"`
}

type RemoveLabelBody struct {
	LabelID uint64 `json:"labelID" binding:"required"`
}

func (h *BoardHandler) RemoveLabel(c *gin.Context) {
	var uri RemoveLabelUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorhandlers.HandleError(c, errorhandlers.NewAPIError(http.StatusBadRequest, "Invalid URI parameters"))
		return
	}

	var body RemoveLabelBody
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

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcReq := &pb_board.RemoveLabelRequest{LabelID: body.LabelID}

	res, err := boardClient.RemoveLabel(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, res.Message, nil)
}

type ArchiveBoardUri struct {
	BoardID uint64 `uri:"boardID" binding:"required"`
}

func (h *BoardHandler) ArchiveBoard(c *gin.Context) {
	var uri ArchiveBoardUri
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

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcReq := &pb_board.ArchiveBoardRequest{}

	res, err := boardClient.ArchiveBoard(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, res.Message, nil)
}

type RestoreBoardUri struct {
	BoardID uint64 `uri:"boardID" binding:"required"`
}

func (h *BoardHandler) RestoreBoard(c *gin.Context) {
	var uri RestoreBoardUri
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

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcReq := &pb_board.RestoreBoardRequest{}

	res, err := boardClient.RestoreBoard(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, res.Message, nil)
}

type DeleteBoardUri struct {
	BoardID uint64 `uri:"boardID" binding:"required"`
}

func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	var uri DeleteBoardUri
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

	md := metadata.Pairs("userID", strconv.FormatUint(userID, 10), "boardID", strconv.FormatUint(uri.BoardID, 10))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	grpcReq := &pb_board.DeleteBoardRequest{}

	res, err := boardClient.DeleteBoard(ctx, grpcReq)
	if err != nil {
		errorhandlers.HandleError(c, err)
		return
	}

	responsehandlers.Success(c, http.StatusOK, res.Message, nil)
}

// Helpers

func (h *BoardHandler) CheckVisibility(ctx context.Context, userID, boardID uint64) error {
	authService, err := h.services.GetAuthClient()
	if err != nil {
		return err
	}

	_, err = authService.CheckBoardVisibility(ctx, &pb_auth.CheckBoardVisibilityRequest{
		UserID:  userID,
		BoardID: boardID,
	})

	return err
}
