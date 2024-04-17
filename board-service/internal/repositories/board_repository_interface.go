package repositories

import (
	internal_models "github.com/sm888sm/halten-backend/board-service/internal/models"
	models "github.com/sm888sm/halten-backend/models"
)

type CreateBoardParams struct {
	Board *models.Board
}

type GetBoardByIDParams struct {
	BoardID uint64
	UserID  uint64
}

type GetBoardListParams struct {
	PageNumber uint64
	PageSize   uint64
	UserID     uint64
}

type GetBoardUsersParams struct {
	BoardID uint64
}

type UpdateBoardParams struct {
	BoardID uint64
	Name    string
}

type AddBoardUsersParams struct {
	BoardID uint64
	UserIDs []uint
	Role    string
}

type RemoveBoardUsersParams struct {
	BoardID uint64
	UserIDs []uint
}

type AssignBoardUserRoleParams struct {
	BoardID      uint64
	AssignUserID uint
	Role         string
}

type ChangeBoardOwnerParams struct {
	BoardID        uint64
	CurrentOwnerID uint64
	NewOwnerID     uint64
}

type GetArchivedBoardListParams struct {
	PageNumber uint64
	PageSize   uint64
	UserID     uint64
}

type RestoreBoardParams struct {
	BoardID uint64
}

type ArchiveBoardParams struct {
	BoardID uint64
}

type DeleteBoardParams struct {
	BoardID uint64
}

type GetBoardPermissionByCardParams struct {
	CardID uint
	UserID uint
}

type GetBoardPermissionByCardResult struct {
	BoardID    uint
	Visibility string
	Permission *models.Permission
}

type GetBoardPermissionByListParams struct {
	ListID uint
	UserID uint
}

type GetBoardPermissionByListResult struct {
	BoardID    uint
	Visibility string
	Permission *models.Permission
}

type BoardRepository interface {
	CreateBoard(params CreateBoardParams) (*models.Board, error)
	GetBoardByID(params GetBoardByIDParams) (*models.Board, error)
	GetBoardList(params GetBoardListParams) (*internal_models.BoardList, error)
	GetBoardUsers(params GetBoardUsersParams) ([]models.User, error)
	UpdateBoard(params UpdateBoardParams) error
	AddBoardUsers(params AddBoardUsersParams) error
	RemoveBoardUsers(params RemoveBoardUsersParams) error
	AssignBoardUserRole(params AssignBoardUserRoleParams) error
	ChangeBoardOwner(params ChangeBoardOwnerParams) error
	GetArchivedBoardList(params GetArchivedBoardListParams) (*internal_models.BoardList, error)
	RestoreBoard(params RestoreBoardParams) error
	ArchiveBoard(params ArchiveBoardParams) error
	DeleteBoard(params DeleteBoardParams) error
	GetBoardPermissionByCard(params GetBoardPermissionByCardParams) (*GetBoardPermissionByCardResult, error)
	GetBoardPermissionByList(params GetBoardPermissionByListParams) (*GetBoardPermissionByListResult, error)
}
