package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	"github.com/sm888sm/halten-backend/common/responsehandler"
	external_services "github.com/sm888sm/halten-backend/gateway-service/external/services"
	pb_list "github.com/sm888sm/halten-backend/list-service/api/pb"
	pb_user "github.com/sm888sm/halten-backend/user-service/api/pb"
)

type CreateListInput struct {
	BoardId uint64 `json:"boardId"`
	Name    string `json:"name"`
}

type UpdateListInput struct {
	Name string `json:"name"`
}

type MoveListPositionInput struct {
	NewPosition uint64 `json:"newPosition"`
}

type ListHandler struct {
	services *external_services.Services
}

func NewListHandler(services *external_services.Services) *ListHandler {
	return &ListHandler{services: services}
}

func (h *ListHandler) CreateList(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).UserID

	var input CreateListInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewHttpBadRequestError())
		return
	}

	listService, err := h.services.GetListClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_list.CreateListRequest{
		UserId:  userId,
		BoardId: input.BoardId,
		Name:    input.Name,
	}
	_, err = listService.CreateList(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "List created successfully", nil)
}

func (h *ListHandler) GetListsByBoard(c *gin.Context) {
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

	listService, err := h.services.GetListClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_list.GetListsByBoardRequest{BoardId: boardID, UserId: uint64(userId)}
	resp, err := listService.GetListsByBoard(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "Lists fetched successfully", resp)
}

func (h *ListHandler) UpdateList(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).UserID

	listIDStr := c.Param("id")
	listID, err := strconv.ParseUint(listIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid list ID"))
		return
	}

	var input UpdateListInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewHttpBadRequestError())
		return
	}

	listService, err := h.services.GetListClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_list.UpdateListRequest{
		Id:     listID,
		UserId: userId,
		Name:   input.Name,
	}

	_, err = listService.UpdateList(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "List updated successfully", nil)
}

func (h *ListHandler) DeleteList(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).UserID

	listIDStr := c.Param("id")
	listID, err := strconv.ParseUint(listIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid list ID"))
		return
	}

	listService, err := h.services.GetListClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_list.DeleteListRequest{
		Id:     listID,
		UserId: userId,
	}

	_, err = listService.DeleteList(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "List deleted successfully", nil)
}

func (h *ListHandler) MoveListPosition(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	userId := user.(*pb_user.User).UserID

	listIDStr := c.Param("id")
	listID, err := strconv.ParseUint(listIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewAPIError(http.StatusBadRequest, "Invalid list ID"))
		return
	}

	var input MoveListPositionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorhandler.NewHttpBadRequestError())
		return
	}

	listService, err := h.services.GetListClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	req := &pb_list.MoveListPositionRequest{
		Id:          listID,
		UserId:      userId,
		NewPosition: input.NewPosition,
	}

	_, err = listService.MoveListPosition(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandler.NewHttpInternalError())
		return
	}

	responsehandler.Success(c, http.StatusOK, "List position moved successfully", nil)
}
