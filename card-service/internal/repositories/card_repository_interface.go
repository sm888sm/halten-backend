package repositories

import (
	"time"

	internal_models "github.com/sm888sm/halten-backend/card-service/internal/models"
	models "github.com/sm888sm/halten-backend/models"
)

type CreateCardRequest struct {
	Card *models.Card
}

type CreateCardResponse struct {
	Card *models.Card
}

type GetCardByIDRequest struct {
	CardID uint64
}

type GetCardByIDResponse struct {
	Card *internal_models.CardDTO
}

type GetCardsByListRequest struct {
	ListID uint64
}

type GetCardsByListResponse struct {
	Cards []*internal_models.CardMetaDTO
}

type GetCardsByBoardRequest struct {
	BoardID uint64
}

type GetCardsByBoardResponse struct {
	Cards []*internal_models.CardMetaDTO
}

type DeleteCardRequest struct {
	CardID  uint64
	BoardID uint64
}

type MoveCardPositionRequest struct {
	CardID      uint64
	NewPosition int64
	BoardID     uint64
	OldListID   uint64
	NewListID   uint64
}

type UpdateCardNameRequest struct {
	CardID  uint64
	Name    string
	BoardID uint64
}

type UpdateCardDescriptionRequest struct {
	CardID      uint64
	Description string
	BoardID     uint64
}

type AddCardLabelRequest struct {
	LabelID uint64
	CardID  uint64
	BoardID uint64
}

type RemoveCardLabelRequest struct {
	LabelID uint64
	CardID  uint64
	BoardID uint64
}

type SetCardDatesRequest struct {
	StartDate *time.Time
	DueDate   *time.Time
	CardID    uint64
	BoardID   uint64
}

type MarkCardCompleteRequest struct {
	CardID  uint64
	BoardID uint64
}

type AddCardAttachmentRequest struct {
	AttachmentID uint64
	CardID       uint64
	BoardID      uint64
}

type RemoveCardAttachmentRequest struct {
	AttachmentID uint64
	CardID       uint64
	BoardID      uint64
}

type AddCardCommentRequest struct {
	Comment models.Comment
	CardID  uint64
	BoardID uint64
	UserID  uint64
}

type RemoveCardCommentRequest struct {
	CommentID uint64
	CardID    uint64
	BoardID   uint64
	UserID    uint64
}

type AddCardMembersRequest struct {
	UserIDs []uint64
	CardID  uint64
	BoardID uint64
}

type RemoveCardMembersRequest struct {
	UserIDs []uint64
	CardID  uint64
	BoardID uint64
}

type ArchiveCardRequest struct {
	CardID  uint64
	BoardID uint64
}

type RestoreCardRequest struct {
	CardID  uint64
	BoardID uint64
}

type CardRepository interface {
	CreateCard(req *CreateCardRequest) (*CreateCardResponse, error)
	GetCardByID(req *GetCardByIDRequest) (*GetCardByIDResponse, error)
	GetCardsByList(req *GetCardsByListRequest) (*GetCardsByListResponse, error)
	GetCardsByBoard(req *GetCardsByBoardRequest) (*GetCardsByBoardResponse, error)
	MoveCardPosition(req *MoveCardPositionRequest) error
	UpdateCardName(req *UpdateCardNameRequest) error
	UpdateCardDescription(req *UpdateCardDescriptionRequest) error
	AddCardLabel(req *AddCardLabelRequest) error
	RemoveCardLabel(req *RemoveCardLabelRequest) error
	SetCardDates(req *SetCardDatesRequest) error
	MarkCardComplete(req *MarkCardCompleteRequest) error
	AddCardAttachment(req *AddCardAttachmentRequest) error
	RemoveCardAttachment(req *RemoveCardAttachmentRequest) error
	AddCardComment(req *AddCardCommentRequest) error
	RemoveCardComment(req *RemoveCardCommentRequest) error
	AddCardMembers(req *AddCardMembersRequest) error
	RemoveCardMembers(req *RemoveCardMembersRequest) error
	ArchiveCard(req *ArchiveCardRequest) error
	RestoreCard(req *RestoreCardRequest) error
	DeleteCard(req *DeleteCardRequest) error
}
