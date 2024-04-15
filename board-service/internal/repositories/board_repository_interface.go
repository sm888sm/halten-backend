package repositories

import (
	internal_models "github.com/sm888sm/halten-backend/board-service/internal/models"
	models "github.com/sm888sm/halten-backend/models"
)

type BoardRepository interface {
	CreateBoard(board *models.Board, userID uint) (*models.Board, error)
	GetBoardByID(id uint, userID uint) (*models.Board, error)
	GetBoardList(userID uint, pageNumber, pageSize int) (*internal_models.BoardList, error)
	GetBoardUsers(boardID uint) ([]models.User, error)
	UpdateBoard(userID uint, id uint, name string) error
	DeleteBoard(userID uint, id uint) error
	AddBoardUsers(userID uint, boardID uint, userIDs []uint, role string) error
	RemoveBoardUsers(userID uint, boardID uint, userIDs []uint) error
	AssignBoardUserRole(userID uint, boardID uint, AssignUserID uint, role string) error
	ChangeBoardOwner(boardID uint, currentOwnerID uint, newOwnerID uint) error
	GetArchivedBoardList(userID uint, pageNumber, pageSize int) (*internal_models.BoardList, error)
	RestoreBoard(userID, boardID uint) error
	ArchiveBoard(userID, boardID uint) error
}
