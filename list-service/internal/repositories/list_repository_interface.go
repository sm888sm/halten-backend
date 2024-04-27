package repositories

import (
	models "github.com/sm888sm/halten-backend/models"
)

type CreateListRequest struct {
	List   *models.List
	UserID uint64
}

type CreateListResponse struct {
	List *models.List
}

type GetListRequest struct {
	ID      uint64
	BoardID uint64
	UserID  uint64
}

type GetListResponse struct {
	List *models.List
}

type GetListsByBoardRequest struct {
	BoardID uint64
	UserID  uint64
}

type GetListsByBoardResponse struct {
	Lists []*models.List
}

type UpdateListNameRequest struct {
	ID      uint64
	Name    string
	BoardID uint64
	UserID  uint64
}

type ArchiveListRequest struct {
	ListID  uint64
	BoardID uint64
}

type RestoreListRequest struct {
	ListID  uint64
	BoardID uint64
}

type DeleteListRequest struct {
	ID      uint64
	BoardID uint64
	UserID  uint64
}

type MoveListPositionRequest struct {
	ID          uint64
	NewPosition int64
	BoardID     uint64
	UserID      uint64
}

type ListRepository interface {
	CreateList(req *CreateListRequest) (*CreateListResponse, error)
	GetListByID(req *GetListRequest) (*GetListResponse, error)
	GetListsByBoard(req *GetListsByBoardRequest) (*GetListsByBoardResponse, error)
	UpdateListName(req *UpdateListNameRequest) error
	MoveListPosition(req *MoveListPositionRequest) error
	ArchiveList(req *ArchiveListRequest) error
	RestoreList(req *RestoreListRequest) error
	DeleteList(req *DeleteListRequest) error
}
