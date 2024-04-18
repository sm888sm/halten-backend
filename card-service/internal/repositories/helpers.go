package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/errorhandler"

	models "github.com/sm888sm/halten-backend/models"
)

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