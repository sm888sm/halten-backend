package repositories

import (
	"time"

	internal_models "github.com/sm888sm/halten-backend/card-service/internal/models"
	models "github.com/sm888sm/halten-backend/models"
)

type CreateCardParams struct {
	Card *models.Card
}

type GetCardByIDParams struct {
	CardID uint64
}

type GetCardsByListParams struct {
	ListID uint64
}

type GetCardsByBoardParams struct {
	BoardID uint64
}

type DeleteCardParams struct {
	CardID  uint64
	BoardID uint64
}

type MoveCardPositionParams struct {
	CardID      uint64
	NewPosition int
	BoardID     uint64
	OldListID   uint64
	NewListID   uint64
}

type UpdateCardNameParams struct {
	CardID  uint64
	Name    string
	BoardID uint64
}

type UpdateCardDescriptionParams struct {
	CardID         uint64
	NewDescription string
	BoardID        uint64
}

type AddCardLabelParams struct {
	LabelID uint64
	CardID  uint64
	BoardID uint64
}

type RemoveCardLabelParams struct {
	LabelID uint64
	CardID  uint64
	BoardID uint64
}

type SetCardDatesParams struct {
	StartDate *time.Time
	DueDate   *time.Time
	CardID    uint64
	BoardID   uint64
}

type MarkCardCompleteParams struct {
	CardID  uint64
	BoardID uint64
}

type AddCardAttachmentParams struct {
	AttachmentID uint64
	CardID       uint64
	BoardID      uint64
}

type RemoveCardAttachmentParams struct {
	AttachmentID uint64
	CardID       uint64
	BoardID      uint64
}

type AddCardCommentParams struct {
	Comment models.Comment
	CardID  uint64
	BoardID uint64
	UserID  uint64
}

type RemoveCardCommentParams struct {
	CommentID uint64
	CardID    uint64
	BoardID   uint64
	UserID    uint64
}

type AddCardMembersParams struct {
	UserIDs []uint64
	CardID  uint64
	BoardID uint64
}

type RemoveCardMembersParams struct {
	UserIDs []uint64
	CardID  uint64
	BoardID uint64
}

type ArchiveCardParams struct {
	CardID  uint64
	BoardID uint64
}

type RestoreCardParams struct {
	CardID  uint64
	BoardID uint64
}

type CardRepository interface {
	CreateCard(params CreateCardParams) error
	GetCardByID(params GetCardByIDParams) (*models.Card, error)
	GetCardsByList(params GetCardsByListParams) ([]*internal_models.CardMeta, error)
	GetCardsByBoard(params GetCardsByBoardParams) ([]*internal_models.CardMeta, error)
	MoveCardPosition(params MoveCardPositionParams) error
	UpdateCardName(params UpdateCardNameParams) error
	UpdateCardDescription(params UpdateCardDescriptionParams) error
	AddCardLabel(params AddCardLabelParams) error
	RemoveCardLabel(params RemoveCardLabelParams) error
	SetCardDates(params SetCardDatesParams) error
	MarkCardComplete(params MarkCardCompleteParams) error
	AddCardAttachment(params AddCardAttachmentParams) error
	RemoveCardAttachment(params RemoveCardAttachmentParams) error
	AddCardComment(params AddCardCommentParams) error
	RemoveCardComment(params RemoveCardCommentParams) error
	AddCardMembers(params AddCardMembersParams) error
	RemoveCardMembers(params RemoveCardMembersParams) error
	ArchiveCard(params ArchiveCardParams) error
	RestoreCard(params RestoreCardParams) error
	DeleteCard(params DeleteCardParams) error
}
