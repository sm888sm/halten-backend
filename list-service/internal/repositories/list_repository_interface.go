package repositories

import (
	models "github.com/sm888sm/halten-backend/models"
)

type ListRepository interface {
	CreateList(list *models.List, userID uint) error
	GetList(id uint, boardID uint, userID uint) (*models.List, error)
	GetListsByBoard(boardID uint, userID uint) ([]*models.List, error)
	UpdateList(id uint, name string, boardID uint, userID uint) error
	DeleteList(id uint, boardID uint, userID uint) error
	MoveListPosition(id uint, newPosition int, boardID uint, userID uint) error
}
