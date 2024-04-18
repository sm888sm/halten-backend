package repositories

import (
	internal_models "github.com/sm888sm/halten-backend/board-service/internal/models"
	models "github.com/sm888sm/halten-backend/models"
)

type CreateBoardRequest struct {
	Board *models.Board
}

type CreateBoardResponse struct {
	Board *models.Board
}

type GetBoardByIDRequest struct {
	BoardID uint64
	UserID  uint64
}

type GetBoardByIDResponse struct {
	Board *models.Board
}

type GetBoardListRequest struct {
	PageNumber uint64
	PageSize   uint64
	UserID     uint64
}

type GetBoardListResponse struct {
	Boards     []*internal_models.BoardMetaDTO
	Pagination *internal_models.Pagination
}

type GetBoardMembersRequest struct {
	BoardID uint64
}

type GetBoardMembersResponse struct {
	Members []*internal_models.MemberDTO
}

type UpdateBoardNameRequest struct {
	BoardID uint64
	Name    string
}

type AddBoardUsersRequest struct {
	BoardID uint64
	UserIDs []uint64
	Role    string
}

type RemoveBoardUsersRequest struct {
	BoardID uint64
	UserIDs []uint64
}

type AssignBoardUserRoleRequest struct {
	BoardID uint64
	UserID  uint64
	UserIDs []uint64
	Role    string
}

type ChangeBoardOwnerRequest struct {
	BoardID        uint64
	CurrentOwnerID uint64
	NewOwnerID     uint64
}

type GetArchivedBoardListRequest struct {
	PageNumber uint64
	PageSize   uint64
	UserID     uint64
}

type GetArchivedBoardListResponse struct {
	Boards     []*internal_models.BoardMetaDTO
	Pagination *internal_models.Pagination
}

type RestoreBoardRequest struct {
	BoardID uint64
}

type ArchiveBoardRequest struct {
	BoardID uint64
}

type DeleteBoardRequest struct {
	BoardID uint64
}

type BoardRepository interface {
	CreateBoard(req *CreateBoardRequest) (*CreateBoardResponse, error)
	GetBoardByID(req *GetBoardByIDRequest) (*GetBoardByIDResponse, error)
	GetBoardList(req *GetBoardListRequest) (*GetBoardListResponse, error)
	GetBoardMembers(req *GetBoardMembersRequest) (*GetBoardMembersResponse, error)
	UpdateBoardName(req *UpdateBoardNameRequest) error
	AddBoardUsers(req *AddBoardUsersRequest) error
	RemoveBoardUsers(req *RemoveBoardUsersRequest) error
	AssignBoardUserRole(req *AssignBoardUserRoleRequest) error
	ChangeBoardOwner(req *ChangeBoardOwnerRequest) error
	GetArchivedBoardList(req *GetArchivedBoardListRequest) (*GetArchivedBoardListResponse, error)
	RestoreBoard(req *RestoreBoardRequest) error
	ArchiveBoard(req *ArchiveBoardRequest) error
	DeleteBoard(req *DeleteBoardRequest) error
}
