package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/sm888sm/halten-backend/common/errorhandlers"

	models "github.com/sm888sm/halten-backend/models"
)

func (r *GormCardRepository) checkLabelExistsAndBelongsToBoard(tx *gorm.DB, labelID uint64, boardID uint64) (*models.Label, error) {
	var existingLabel models.Label
	if err := tx.First(&existingLabel, labelID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandlers.NewGrpcNotFoundError("Label not found")
		}
		return nil, errorhandlers.NewGrpcInternalError()
	}

	if existingLabel.BoardID != boardID {
		return nil, errorhandlers.NewGrpcNotFoundError("Label does not belong to the board")
	}

	return &existingLabel, nil
}

func (r *GormCardRepository) checkCardExistsAndBelongsToBoard(tx *gorm.DB, cardID uint64, boardID uint64) (*models.Card, error) {
	card := &models.Card{BaseModel: models.BaseModel{ID: cardID}, BoardID: boardID}
	if err := tx.First(card).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandlers.NewGrpcNotFoundError("Card not found")
		}
		return nil, errorhandlers.NewGrpcInternalError()
	}

	return card, nil
}
