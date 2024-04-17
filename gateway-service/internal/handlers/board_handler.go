package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	pb_board "github.com/sm888sm/halten-backend/board-service/api/pb"
	pb_user "github.com/sm888sm/halten-backend/user-service/api/pb"
	"google.golang.org/grpc/metadata"

	"github.com/gin-gonic/gin"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	"github.com/sm888sm/halten-backend/common/responsehandler"
	external_services "github.com/sm888sm/halten-backend/gateway-service/external/services"
)

type CreateBoardInput struct {
	Name string `json:"name"`
}

type UpdateBoardInput struct {
	Name string `json:"name"`
}

type AddBoardUsersInput struct {
	UserIds []uint64 `json:"userIds"`
}

type RemoveBoardUsersInput struct {
	UserIds []uint64 `json:"userIds"`
}

type ChangeBoardOwnerInput struct {
	NewOwnerId uint64 `json:"newOwnerId"`
}

type BoardHandler struct {
	services *external_services.Services
}

type AssignUserRoleInput struct {
	UserId uint64 `json:"userId"`
	Role   string `json:"role"`
}

func NewBoardHandler(services *external_services.Services) *BoardHandler {
	return &BoardHandler{services: services}
}

func (h *BoardHandler) CreateBoard(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var input CreateBoardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewHttpBadRequestError())
		return
	}

	md := metadata.Pairs("userID", fmt.Sprintf("%d", userID))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	boardService, err := h.services.GetBoardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_board.CreateBoardRequest{
		UserID: userID,
		Name:   input.Name,
	}
	_, err = boardService.CreateBoard(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Board created successfully", nil)
}

func (h *BoardHandler) GetBoardByID(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	boardIDStr := c.Param("board-id")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid board ID"))
		return
	}

	boardService, err := h.services.GetBoardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_board.GetBoardByIDRequest{BoardID: boardID, UserID: userID}
	resp, err := boardService.GetBoardByID(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Board fetched successfully", resp)
}

func (h *BoardHandler) GetBoardList(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	pageNumberStr, pageSizeStr := c.DefaultQuery("page_number", "1"), c.DefaultQuery("page_size", "10")

	pageNumber, err := strconv.ParseUint(pageNumberStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid page number"))
		return
	}

	pageSize, err := strconv.ParseUint(pageSizeStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid page size"))
		return
	}

	md := metadata.Pairs("userID", fmt.Sprintf("%d", userID))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	boardService, err := h.services.GetBoardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_board.GetBoardListRequest{
		UserID:     userID,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
	resp, err := boardService.GetBoardList(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	pagination := responsehandler.Pagination{
		CurrentPage:  int(resp.Pagination.CurrentPage),
		TotalPages:   int(resp.Pagination.TotalItems),
		ItemsPerPage: int(resp.Pagination.ItemsPerPage),
		TotalItems:   int(resp.Pagination.TotalItems),
		HasMore:      resp.Pagination.HasMore,
	}

	responsehandler.SuccessWithPagination(c, http.StatusOK, "Boards fetched successfully", resp, &pagination)
}

func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).UserID

	boardIDStr := c.Param("id")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid board ID"))
		return
	}

	var input UpdateBoardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewHttpBadRequestError())
		return
	}

	boardService, err := h.services.GetBoardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_board.UpdateBoardRequest{
		Id:     boardID,
		UserId: userId,
		Name:   input.Name,
	}

	_, err = boardService.UpdateBoard(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Board updated successfully", nil)
}

func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).UserID

	boardIDStr := c.Param("id")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid board ID"))
		return
	}

	boardService, err := h.services.GetBoardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_board.DeleteBoardRequest{
		Id:     boardID,
		UserId: userId,
	}
	_, err = boardService.DeleteBoard(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Board deleted successfully", nil)
}

func (h *BoardHandler) AddBoardUsers(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).UserID

	boardIDStr := c.Param("id")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid board ID"))
		return
	}

	var input AddBoardUsersInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid user ID list"))
		return
	}

	boardService, err := h.services.GetBoardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_board.AddUsersRequest{
		Id:      boardID,
		UserId:  userId,
		UserIds: input.UserIds,
	}
	_, err = boardService.AddBoardUser(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Users added to board successfully", nil)
}

func (h *BoardHandler) RemoveBoardUsers(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).UserID

	boardIDStr := c.Param("id")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid board ID"))
		return
	}

	var input RemoveBoardUsersInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewHttpBadRequestError())
		return
	}

	boardService, err := h.services.GetBoardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_board.RemoveUsersRequest{
		Id:      boardID,
		UserId:  userId,
		UserIds: input.UserIds,
	}
	_, err = boardService.RemoveBoardUser(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Users removed from board successfully", nil)
}

func (h *BoardHandler) GetBoardUsers(c *gin.Context) {
	boardIDStr := c.Param("id")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid board ID"))
		return
	}

	boardService, err := h.services.GetBoardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_board.GetBoardUsersRequest{
		Id: boardID,
	}
	resp, err := boardService.GetBoardUsers(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Board users retrieved successfully", resp.Users)
}

func (h *BoardHandler) AssignBoardUserRole(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).UserID

	boardIDStr := c.Param("id")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid board ID"))
		return
	}

	var input AssignUserRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewHttpBadRequestError())
		return
	}

	boardService, err := h.services.GetBoardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_board.AssignUserRoleRequest{
		Id:           boardID,
		UserId:       userId,
		AssignUserId: input.UserId,
		Role:         input.Role,
	}
	_, err = boardService.AssignBoardUserRole(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "User role assigned successfully", nil)
}

func (h *BoardHandler) ChangeBoardOwner(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).UserID

	boardIDStr := c.Param("id")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid board ID"))
		return
	}

	var input ChangeBoardOwnerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewHttpBadRequestError())
		return
	}

	boardService, err := h.services.GetBoardClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_board.ChangeBoardOwnerRequest{
		Id:         boardID,
		UserId:     userId,
		NewOwnerId: input.NewOwnerId,
	}
	_, err = boardService.ChangeBoardOwner(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Board owner changed successfully", nil)
}
