package repositories

import (
	"errors"
	"sort"

	"gorm.io/gorm"

	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/errorhandler"

	internal_models "github.com/sm888sm/halten-backend/card-service/internal/models"
	models "github.com/sm888sm/halten-backend/models"
)

type GormCardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *GormCardRepository {
	return &GormCardRepository{db: db}
}

func (r *GormCardRepository) CreateCard(params CreateCardParams) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(params.Card).Error; err != nil {
			return errorhandler.NewAPIError(httpcodes.ErrBadRequest, err.Error())
		}
		return nil
	})
}

func (r *GormCardRepository) GetCardByID(params GetCardByIDParams) (*models.Card, error) {
	var card models.Card
	if err := r.db.Preload("Attachments", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}).Preload("Labels", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}).Preload("Members", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}).Select("ID", "BoardID", "ListID", "Name", "Description", "Position", "IsCompleted", "StartDate", "DueDate", "CreatedAt", "UpdatedAt").
		Where("id = ? AND archived = false", params.CardID).
		First(&card).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NewAPIError(httpcodes.ErrNotFound, "Card not found")
		}
		return nil, errorhandler.NewGrpcInternalError()
	}

	return &card, nil
}

func (r *GormCardRepository) GetCardsByList(params GetCardsByListParams) ([]*internal_models.CardMeta, error) {
	var cards []*internal_models.CardMeta
	if err := r.db.Model(&models.Card{}).
		Select("id, list_id, board_id, name, position, is_completed, start_date, due_date, created_at, updated_at, "+
			"(SELECT COUNT(*) FROM attachments WHERE card_id = cards.id) as total_attachment, "+
			"(SELECT COUNT(*) FROM comments WHERE card_id = cards.id) as total_comment").
		Preload("Labels", func(db *gorm.DB) *gorm.DB {
			return db.Select("id")
		}).
		Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Select("id")
		}).
		Where("list_id = ? AND archived = false", params.ListID).
		Find(&cards).Error; err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	return cards, nil
}

func (r *GormCardRepository) GetCardsByBoard(params GetCardsByBoardParams) ([]*internal_models.CardMeta, error) {
	var cards []*internal_models.CardMeta
	if err := r.db.Model(&models.Card{}).
		Select("id, list_id, board_id, name, position, is_completed, start_date, due_date, created_at, updated_at, "+
			"(SELECT COUNT(*) FROM attachments WHERE card_id = cards.id) as total_attachment, "+
			"(SELECT COUNT(*) FROM comments WHERE card_id = cards.id) as total_comment").
		Preload("Labels", func(db *gorm.DB) *gorm.DB {
			return db.Select("id")
		}).
		Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Select("id")
		}).
		Where("board_id = ?", params.BoardID).
		Find(&cards).Error; err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	return cards, nil
}

func (r *GormCardRepository) MoveCardPosition(params MoveCardPositionParams) error {
	var count int64
	r.db.Model(&models.Card{}).Where("id = ? AND list_id = ? AND board_id = ?", params.CardID, params.OldListID, params.BoardID).Count(&count)
	if count == 0 {
		return errorhandler.NewAPIError(httpcodes.ErrNotFound, "Card not found")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get all cards in the old list
		var oldCards []*models.Card
		if err := tx.Where("list_id = ?", params.OldListID).Find(&oldCards).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		// Sort cards by position
		sort.Slice(oldCards, func(i, j int) bool {
			return oldCards[i].Position < oldCards[j].Position
		})

		// Find the card to be moved
		var movingCard *models.Card
		for i, c := range oldCards {
			if c.ID == uint(params.CardID) {
				movingCard = c
				oldCards = append(oldCards[:i], oldCards[i+1:]...)
				break
			}
		}

		// Update the positions of the remaining cards in the old list
		for i, c := range oldCards {
			c.Position = i + 1
			if err := tx.Save(&c).Error; err != nil {
				return err
			}
		}

		// If the card is moving to a different list, update the positions in the new list
		if params.OldListID != params.NewListID {
			var newCards []*models.Card
			if err := tx.Where("list_id = ?", params.NewListID).Find(&newCards).Error; err != nil {
				return errorhandler.NewGrpcInternalError()
			}

			// Sort cards by position
			sort.Slice(newCards, func(i, j int) bool {
				return newCards[i].Position < newCards[j].Position
			})

			// Insert the card at the new position and update the positions of the other cards
			newCards = append(newCards, nil)
			copy(newCards[params.NewPosition+1:], newCards[params.NewPosition:])
			newCards[params.NewPosition] = movingCard
			for i, c := range newCards {
				c.Position = i + 1
				if err := tx.Save(&c).Error; err != nil {
					return err
				}
			}
		} else { // If the card is moving within the same list, just update its position
			movingCard.Position = params.NewPosition
			if err := tx.Save(&movingCard).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *GormCardRepository) UpdateCardName(params UpdateCardNameParams) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		card, err := r.checkCardExistsAndBelongsToBoard(tx, params.CardID, params.BoardID)
		if err != nil {
			return err
		}

		if card.Name != params.Name {
			db := tx.Model(card).Update("Name", params.Name)
			if db.Error != nil {
				return errorhandler.NewGrpcInternalError()
			}

			if db.RowsAffected == 0 {
				return errorhandler.NewAPIError(httpcodes.ErrNotFound, "No card found to update")
			}
		}

		return nil
	})
}

func (r *GormCardRepository) UpdateCardDescription(params UpdateCardDescriptionParams) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		card := &models.Card{Model: gorm.Model{ID: uint(params.CardID)}, BoardID: uint(params.BoardID)}
		if err := tx.First(card).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(httpcodes.ErrNotFound, "Card not found")
			}
			return errorhandler.NewGrpcInternalError()
		}

		if card.Description != params.NewDescription {
			db := tx.Model(card).Update("Description", params.NewDescription)
			if db.Error != nil {
				return errorhandler.NewGrpcInternalError()
			}

			if db.RowsAffected == 0 {
				return errorhandler.NewAPIError(httpcodes.ErrNotFound, "No card found to update")
			}
		}

		return nil
	})
}

func (r *GormCardRepository) AddCardLabel(params AddCardLabelParams) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		existingLabel, err := r.checkLabelExistsAndBelongsToBoard(tx, params.LabelID, params.BoardID)
		if err != nil {
			return err
		}

		if err := tx.Model(&models.Card{Model: gorm.Model{ID: uint(params.CardID)}}).Association("Labels").Append(existingLabel); err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormCardRepository) RemoveCardLabel(params RemoveCardLabelParams) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		card, err := r.checkCardExistsAndBelongsToBoard(tx, params.CardID, params.BoardID)
		if err != nil {
			return err
		}

		existingLabel, err := r.checkLabelExistsAndBelongsToBoard(tx, params.LabelID, params.BoardID)
		if err != nil {
			return err
		}

		var label models.Label
		if err := tx.Model(&card).Association("Labels").Find(&label, "id = ?", params.LabelID).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		if label.ID == 0 {
			return errorhandler.NewAPIError(httpcodes.ErrNotFound, "Label not found in the card")
		}

		if err := tx.Model(&models.Card{Model: gorm.Model{ID: uint(params.CardID)}}).Association("Labels").Delete(existingLabel); err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormCardRepository) SetCardDates(params SetCardDatesParams) error {
	// Ensure startDate is no later than dueDate
	if params.StartDate != nil && params.DueDate != nil && params.StartDate.After(*params.DueDate) {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Start date cannot be later than due date")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		card, err := r.checkCardExistsAndBelongsToBoard(tx, params.CardID, params.BoardID)
		if err != nil {
			return err
		}

		changes := false

		if card.StartDate != params.StartDate {
			card.StartDate = params.StartDate
			changes = true
		}

		if card.DueDate != params.DueDate {
			card.DueDate = params.DueDate
			changes = true
		}

		// If both startDate and dueDate are unset, unmark the card as complete
		if params.StartDate == nil && params.DueDate == nil && card.IsCompleted {
			card.IsCompleted = false
			changes = true
		}

		if changes {
			if err := tx.Save(card).Error; err != nil {
				return errorhandler.NewGrpcInternalError()
			}
		}

		return nil
	})
}

func (r *GormCardRepository) MarkCardComplete(params MarkCardCompleteParams) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		card, err := r.checkCardExistsAndBelongsToBoard(tx, params.CardID, params.BoardID)
		if err != nil {
			return err
		}

		// Only cards with a due date can be marked as complete
		if card.DueDate == nil {
			return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Card cannot be marked as complete without a due date")
		}

		card.IsCompleted = true

		if err := tx.Save(card).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormCardRepository) AddCardAttachment(params AddCardAttachmentParams) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		card, err := r.checkCardExistsAndBelongsToBoard(tx, params.CardID, params.BoardID)
		if err != nil {
			return err
		}

		// Check the number of attachments for the card
		var count int64
		if err := tx.Model(&models.Attachment{}).Where("card_id = ?", params.CardID).Count(&count).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}
		if count >= 10 {
			return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Card cannot have more than 10 attachments")
		}

		attachment := &models.Attachment{Model: gorm.Model{ID: uint(params.AttachmentID)}}
		if err := tx.First(attachment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(httpcodes.ErrNotFound, "Attachment not found")
			}
			return errorhandler.NewGrpcInternalError()
		}

		// Ensure the attachment belongs to the same board
		if attachment.BoardID != card.BoardID {
			return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Attachment does not belong to the same board")
		}

		attachment.CardID = uint(params.CardID)

		if err := tx.Save(attachment).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormCardRepository) RemoveCardAttachment(params RemoveCardAttachmentParams) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if the card exists
		card, err := r.checkCardExistsAndBelongsToBoard(tx, params.CardID, params.BoardID)
		if err != nil {
			return err
		}

		// Check if the attachment exists
		var attachment models.Attachment
		if err := tx.Where("id = ?", params.AttachmentID).First(&attachment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(httpcodes.ErrNotFound, "Attachment not found")
			}
			return errorhandler.NewGrpcInternalError()
		}

		// Remove attachment from card
		if err := tx.Model(&card).Association("Attachments").Delete(&attachment).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormCardRepository) AddCardComment(params AddCardCommentParams) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		_, err := r.checkCardExistsAndBelongsToBoard(tx, params.CardID, params.BoardID)
		if err != nil {
			return err
		}

		params.Comment.CardID = uint(params.CardID)
		params.Comment.UserID = uint(params.UserID)
		if err := tx.Create(&params.Comment).Error; err != nil {
			return errorhandler.NewAPIError(httpcodes.ErrBadRequest, err.Error())
		}
		return nil
	})
}

func (r *GormCardRepository) RemoveCardComment(params RemoveCardCommentParams) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		_, err := r.checkCardExistsAndBelongsToBoard(tx, params.CardID, params.BoardID)
		if err != nil {
			return err
		}

		comment := &models.Comment{Model: gorm.Model{ID: uint(params.CommentID)}, CardID: uint(params.CardID)}
		if err := tx.First(comment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(httpcodes.ErrNotFound, "Comment not found")
			}
			return errorhandler.NewGrpcInternalError()
		}

		// // Check if the user is the admin, the owner of the card, or the one who created the comment
		// if params.UserID != uint64(comment.UserID) {
		// 	return errorhandler.NewAPIError(httpcodes.ErrUnauthorized, "You are not authorized to delete this comment")
		// }

		if err := tx.Delete(comment).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormCardRepository) AddCardMembers(params AddCardMembersParams) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if the card exists
		card, err := r.checkCardExistsAndBelongsToBoard(tx, params.CardID, params.BoardID)
		if err != nil {
			return err
		}

		// For each userID, find the user and add them to the card
		for _, userID := range params.UserIDs {
			var user models.User
			if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errorhandler.NewAPIError(httpcodes.ErrNotFound, "User not found")
				}
				return errorhandler.NewGrpcInternalError()
			}

			// Add user to card
			tx.Model(&card).Association("Members").Append(&user)
		}

		return nil
	})
}

func (r *GormCardRepository) RemoveCardMembers(params RemoveCardMembersParams) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if the card exists
		card, err := r.checkCardExistsAndBelongsToBoard(tx, params.CardID, params.BoardID)
		if err != nil {
			return err
		}

		// For each userID, find the user and remove them from the card
		for _, userID := range params.UserIDs {
			var user models.User
			if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errorhandler.NewAPIError(httpcodes.ErrNotFound, "User not found")
				}
				return errorhandler.NewGrpcInternalError()
			}

			// Remove user from card
			tx.Model(&card).Association("Members").Delete(&user)
		}

		return nil
	})
}

func (r *GormCardRepository) ArchiveCard(params ArchiveCardParams) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		card, err := r.checkCardExistsAndBelongsToBoard(tx, params.CardID, params.BoardID)
		if err != nil {
			return err
		}

		card.IsArchived = true
		if err := tx.Save(card).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormCardRepository) RestoreCard(params RestoreCardParams) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		card, err := r.checkCardExistsAndBelongsToBoard(tx, params.CardID, params.BoardID)
		if err != nil {
			return err
		}

		card.IsArchived = false
		if err := tx.Save(card).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormCardRepository) DeleteCard(params DeleteCardParams) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND board_id = ?", params.CardID, params.BoardID).Delete(&models.Card{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(httpcodes.ErrNotFound, "Card not found")
			}
			return errorhandler.NewGrpcInternalError()
		}
		return nil
	})
}

// Helpers

func (r *GormCardRepository) checkLabelExistsAndBelongsToBoard(tx *gorm.DB, labelID uint64, boardID uint64) (*models.Label, error) {
	var existingLabel models.Label
	if err := tx.First(&existingLabel, labelID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NewAPIError(httpcodes.ErrNotFound, "Label not found")
		}
		return nil, errorhandler.NewGrpcInternalError()
	}

	if existingLabel.BoardID != uint(boardID) {
		return nil, errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Label does not belong to the board")
	}

	return &existingLabel, nil
}

func (r *GormCardRepository) checkCardExistsAndBelongsToBoard(tx *gorm.DB, cardID uint64, boardID uint64) (*models.Card, error) {
	card := &models.Card{Model: gorm.Model{ID: uint(cardID)}, BoardID: uint(boardID)}
	if err := tx.First(card).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NewAPIError(httpcodes.ErrNotFound, "Card not found")
		}
		return nil, errorhandler.NewGrpcInternalError()
	}

	return card, nil
}
