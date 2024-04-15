package repositories

import (
	"time"

	internal_models "github.com/sm888sm/halten-backend/card-service/internal/models"
	models "github.com/sm888sm/halten-backend/models"
)

type CardRepository interface {
	CreateCard(card *models.Card, userID uint) error
	GetCardByID(cardID uint, boardID uint, userID uint) (*models.Card, error)
	GetCardsByList(listID uint, boardID uint, userID uint) ([]*internal_models.CardMeta, error)
	GetCardsByBoard(boardID uint, userID uint) ([]*internal_models.CardMeta, error)
	DeleteCard(cardID uint, boardID uint, userID uint) error
	MoveCardPosition(cardID uint, newPosition int, boardID uint, oldListID uint, newListID uint, userID uint) error
	UpdateCardName(cardID uint, newName string, boardID uint, userID uint) error
	UpdateCardDescription(cardID uint, newDescription string, boardID uint, userID uint) error
	AddCardLabel(label models.Label, cardID uint, boardID uint, userID uint) error
	RemoveCardLabel(labelID uint, cardID uint, boardID uint, userID uint) error
	SetCardDates(startDate, dueDate *time.Time, cardID uint, boardID uint, userID uint) error
	MarkCardComplete(cardID uint, boardID uint, userID uint) error
	AddCardAttachment(attachmentID uint, cardID uint, boardID uint, userID uint) error
	RemoveCardAttachment(attachmentID uint, cardID uint, boardID uint, userID uint) error
	AddCardComment(comment models.Comment, cardID uint, boardID uint, userID uint) error
	RemoveCardComment(commentID uint, cardID uint, boardID uint, userID uint) error
	AddCardMembers(userIDs []uint, cardID uint, boardID uint, userID uint) error
	RemoveCardMembers(userIDs []uint, cardID uint, boardID uint, userID uint) error
	ArchiveCard(cardID uint, boardID uint, userID uint) error
	RestoreCard(cardID uint, boardID uint, userID uint) error
}
