package repositories

import (
	"errors"
	"sort"
	"time"

	internal_models "github.com/sm888sm/halten-backend/card-service/internal/models"
	"github.com/sm888sm/halten-backend/common"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	models "github.com/sm888sm/halten-backend/models"

	"gorm.io/gorm"
)

type GormCardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *GormCardRepository {
	return &GormCardRepository{db: db}
}

func (r *GormCardRepository) CreateCard(card *models.Card, userID uint) error {
	if err := r.checkPermission(card.BoardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(card).Error; err != nil {
			return errorhandler.NewAPIError(errorhandler.ErrBadRequest, err.Error())
		}
		return nil
	})
}

func (r *GormCardRepository) GetCardByID(cardID uint, boardID uint, userID uint) (*models.Card, error) {
	if err := r.checkPermission(boardID, userID); err != nil {
		return nil, err
	}

	var card models.Card
	if err := r.db.Preload("Attachments").Preload("Labels").Preload("Members").
		Select("ID", "BoardID", "ListID", "Name", "Description", "Position", "IsCompleted", "StartDate", "DueDate", "CreatedAt", "UpdatedAt").
		Where("id = ? AND board_id = ? AND archived = false", cardID, boardID).
		First(&card).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NewAPIError(errorhandler.ErrNotFound, "Card not found")
		}
		return nil, errorhandler.NewGrpcInternalError()
	}

	return &card, nil
}

func (r *GormCardRepository) GetCardsByList(listID uint, boardID uint, userID uint) ([]*internal_models.CardMeta, error) {
	if err := r.checkPermission(boardID, userID); err != nil {
		return nil, err
	}

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
		Where("list_id = ? AND board_id = ? AND archived = false", listID, boardID).
		Find(&cards).Error; err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	return cards, nil
}

func (r *GormCardRepository) GetCardsByBoard(boardID uint, userID uint) ([]*internal_models.CardMeta, error) {
	if err := r.checkPermission(boardID, userID); err != nil {
		return nil, err
	}

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
		Where("board_id = ?", boardID).
		Find(&cards).Error; err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	return cards, nil
}

func (r *GormCardRepository) DeleteCard(cardID uint, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND board_id = ?", cardID, boardID).Delete(&models.Card{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(errorhandler.ErrNotFound, "Card not found")
			}
			return errorhandler.NewGrpcInternalError()
		}
		return nil
	})
}

func (r *GormCardRepository) MoveCardPosition(cardID uint, newPosition int, boardID uint, oldListID uint, newListID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	var count int64
	r.db.Model(&models.Card{}).Where("id = ? AND list_id = ? AND board_id = ?", cardID, oldListID, boardID).Count(&count)
	if count == 0 {
		return errorhandler.NewAPIError(errorhandler.ErrNotFound, "Card not found")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get all cards in the old list
		var oldCards []*models.Card
		if err := tx.Where("list_id = ?", oldListID).Find(&oldCards).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		// Sort cards by position
		sort.Slice(oldCards, func(i, j int) bool {
			return oldCards[i].Position < oldCards[j].Position
		})

		// Find the card to be moved
		var movingCard *models.Card
		for i, c := range oldCards {
			if c.ID == cardID {
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
		if oldListID != newListID {
			var newCards []*models.Card
			if err := tx.Where("list_id = ?", newListID).Find(&newCards).Error; err != nil {
				return errorhandler.NewGrpcInternalError()
			}

			// Sort cards by position
			sort.Slice(newCards, func(i, j int) bool {
				return newCards[i].Position < newCards[j].Position
			})

			// Insert the card at the new position and update the positions of the other cards
			newCards = append(newCards, nil)
			copy(newCards[newPosition+1:], newCards[newPosition:])
			newCards[newPosition] = movingCard
			for i, c := range newCards {
				c.Position = i + 1
				if err := tx.Save(&c).Error; err != nil {
					return err
				}
			}
		} else { // If the card is moving within the same list, just update its position
			movingCard.Position = newPosition
			if err := tx.Save(&movingCard).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *GormCardRepository) UpdateCardName(cardID uint, newName string, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		card, err := r.checkCardExistsAndBelongsToBoard(tx, cardID, boardID)
		if err != nil {
			return err
		}

		if card.Name != newName {
			db := tx.Model(card).Update("Name", newName)
			if db.Error != nil {
				return errorhandler.NewGrpcInternalError()
			}

			if db.RowsAffected == 0 {
				return errorhandler.NewAPIError(errorhandler.ErrNotFound, "No card found to update")
			}
		}

		return nil
	})
}

func (r *GormCardRepository) UpdateCardDescription(cardID uint, newDescription string, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		card := &models.Card{Model: gorm.Model{ID: cardID}, BoardID: boardID}
		if err := tx.First(card).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(errorhandler.ErrNotFound, "Card not found")
			}
			return errorhandler.NewGrpcInternalError()
		}

		if card.Description != newDescription {
			db := tx.Model(card).Update("Description", newDescription)
			if db.Error != nil {
				return errorhandler.NewGrpcInternalError()
			}

			if db.RowsAffected == 0 {
				return errorhandler.NewAPIError(errorhandler.ErrNotFound, "No card found to update")
			}
		}

		return nil
	})
}

func (r *GormCardRepository) AddCardLabel(label models.Label, cardID uint, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		existingLabel, err := r.checkLabelExistsAndBelongsToBoard(tx, label.ID, boardID)
		if err != nil {
			return err
		}

		if err := tx.Model(&models.Card{Model: gorm.Model{ID: cardID}}).Association("Labels").Append(existingLabel); err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormCardRepository) RemoveCardLabel(labelID uint, cardID uint, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		card, err := r.checkCardExistsAndBelongsToBoard(tx, cardID, boardID)
		if err != nil {
			return err
		}

		existingLabel, err := r.checkLabelExistsAndBelongsToBoard(tx, labelID, boardID)
		if err != nil {
			return err
		}

		var label models.Label
		if err := tx.Model(&card).Association("Labels").Find(&label, "id = ?", labelID).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		if label.ID == 0 {
			return errorhandler.NewAPIError(errorhandler.ErrNotFound, "Label not found in the card")
		}

		if err := tx.Model(&models.Card{Model: gorm.Model{ID: cardID}}).Association("Labels").Delete(existingLabel); err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormCardRepository) SetCardDates(startDate, dueDate *time.Time, cardID uint, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	// Ensure startDate is no later than dueDate
	if startDate != nil && dueDate != nil && startDate.After(*dueDate) {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Start date cannot be later than due date")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		card, err := r.checkCardExistsAndBelongsToBoard(tx, cardID, boardID)
		if err != nil {
			return err
		}

		changes := false

		if card.StartDate != startDate {
			card.StartDate = startDate
			changes = true
		}

		if card.DueDate != dueDate {
			card.DueDate = dueDate
			changes = true
		}

		// If both startDate and dueDate are unset, unmark the card as complete
		if startDate == nil && dueDate == nil && card.IsCompleted {
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

func (r *GormCardRepository) MarkCardComplete(cardID uint, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		card, err := r.checkCardExistsAndBelongsToBoard(tx, cardID, boardID)
		if err != nil {
			return err
		}

		// Only cards with a due date can be marked as complete
		if card.DueDate == nil {
			return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Card cannot be marked as complete without a due date")
		}

		card.IsCompleted = true

		if err := tx.Save(card).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormCardRepository) AddCardAttachment(attachmentID uint, cardID uint, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		card, err := r.checkCardExistsAndBelongsToBoard(tx, cardID, boardID)
		if err != nil {
			return err
		}

		// Check the number of attachments for the card
		var count int64
		if err := tx.Model(&models.Attachment{}).Where("card_id = ?", cardID).Count(&count).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}
		if count >= 50 {
			return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Card cannot have more than 50 attachments")
		}

		attachment := &models.Attachment{Model: gorm.Model{ID: attachmentID}}
		if err := tx.First(attachment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(errorhandler.ErrNotFound, "Attachment not found")
			}
			return errorhandler.NewGrpcInternalError()
		}

		// Ensure the attachment belongs to the same board
		if attachment.BoardID != card.BoardID {
			return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Attachment does not belong to the same board")
		}

		attachment.CardID = cardID

		if err := tx.Save(attachment).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormCardRepository) RemoveCardAttachment(attachmentID uint, cardID uint, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if the card exists
		card, err := r.checkCardExistsAndBelongsToBoard(tx, cardID, boardID)
		if err != nil {
			return err
		}

		// Check if the attachment exists
		var attachment models.Attachment
		if err := tx.Where("id = ?", attachmentID).First(&attachment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(errorhandler.ErrNotFound, "Attachment not found")
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

func (r *GormCardRepository) AddCardComment(comment models.Comment, cardID uint, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		_, err := r.checkCardExistsAndBelongsToBoard(tx, cardID, boardID)
		if err != nil {
			return err
		}

		comment.CardID = cardID
		comment.UserID = userID
		if err := tx.Create(&comment).Error; err != nil {
			return errorhandler.NewAPIError(errorhandler.ErrBadRequest, err.Error())
		}
		return nil
	})
}

func (r *GormCardRepository) RemoveCardComment(commentID uint, cardID uint, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		_, err := r.checkCardExistsAndBelongsToBoard(tx, cardID, boardID)
		if err != nil {
			return err
		}

		comment := &models.Comment{Model: gorm.Model{ID: commentID}, CardID: cardID}
		if err := tx.First(comment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(errorhandler.ErrNotFound, "Comment not found")
			}
			return errorhandler.NewGrpcInternalError()
		}

		if err := tx.Delete(comment).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormCardRepository) AddCardMembers(userIDs []uint, cardID uint, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if the card exists
		card, err := r.checkCardExistsAndBelongsToBoard(tx, cardID, boardID)
		if err != nil {
			return err
		}

		// For each userID, find the user and add them to the card
		for _, userID := range userIDs {
			var user models.User
			if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errorhandler.NewAPIError(errorhandler.ErrNotFound, "User not found")
				}
				return errorhandler.NewGrpcInternalError()
			}

			// Add user to card
			tx.Model(&card).Association("Members").Append(&user)
		}

		return nil
	})
}

func (r *GormCardRepository) RemoveCardMembers(userIDs []uint, cardID uint, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if the card exists
		card, err := r.checkCardExistsAndBelongsToBoard(tx, cardID, boardID)
		if err != nil {
			return err
		}

		// For each userID, find the user and remove them from the card
		for _, userID := range userIDs {
			var user models.User
			if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errorhandler.NewAPIError(errorhandler.ErrNotFound, "User not found")
				}
				return errorhandler.NewGrpcInternalError()
			}

			// Remove user from card
			tx.Model(&card).Association("Members").Delete(&user)
		}

		return nil
	})
}

func (r *GormCardRepository) ArchiveCard(cardID uint, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		card, err := r.checkCardExistsAndBelongsToBoard(tx, cardID, boardID)
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

func (r *GormCardRepository) RestoreCard(cardID uint, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		card, err := r.checkCardExistsAndBelongsToBoard(tx, cardID, boardID)
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

// Helpers

func (r *GormCardRepository) checkPermission(boardID uint, userID uint) error {
	var permission models.Permission
	if err := r.db.Where("board_id = ? AND user_id = ?", boardID, userID).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorhandler.NewAPIError(errorhandler.ErrForbidden, "Permission not found")
		}
		return err
	}

	if permission.Role == common.OwnerRole || permission.Role == common.AdminRole || permission.Role == common.MemberRole {
		return nil
	}

	return errorhandler.NewAPIError(errorhandler.ErrForbidden, "User does not have permission to perform this operation")
}

func (r *GormCardRepository) checkLabelExistsAndBelongsToBoard(tx *gorm.DB, labelID uint, boardID uint) (*models.Label, error) {
	var existingLabel models.Label
	if err := tx.First(&existingLabel, labelID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NewAPIError(errorhandler.ErrNotFound, "Label not found")
		}
		return nil, errorhandler.NewGrpcInternalError()
	}

	if existingLabel.BoardID != boardID {
		return nil, errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Label does not belong to the board")
	}

	return &existingLabel, nil
}

func (r *GormCardRepository) checkCardExistsAndBelongsToBoard(tx *gorm.DB, cardID uint, boardID uint) (*models.Card, error) {
	card := &models.Card{Model: gorm.Model{ID: cardID}, BoardID: boardID}
	if err := tx.First(card).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NewAPIError(errorhandler.ErrNotFound, "Card not found")
		}
		return nil, errorhandler.NewGrpcInternalError()
	}

	return card, nil
}
