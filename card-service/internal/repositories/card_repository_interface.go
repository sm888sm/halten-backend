package repositories

import (
	models "github.com/sm888sm/halten-backend/models"
)

type CardRepository interface {
	CreateCard(card *models.Card, userID uint) error
	GetCard(id uint, listID uint, userID uint) (*models.Card, error)
	GetCardsByList(listID uint, userID uint) ([]*models.Card, error)
	UpdateCard(id uint, name string, listID uint, userID uint) error
	DeleteCard(id uint, listID uint, userID uint) error
	MoveCardPosition(id uint, newPosition int, listID uint, userID uint) error
}
