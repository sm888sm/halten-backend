package repositories

import (
	"errors"
	"sort"

	"github.com/sm888sm/halten-backend/common"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	models "github.com/sm888sm/halten-backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormCardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *GormCardRepository {
	return &GormCardRepository{db: db}
}

func (r *GormCardRepository) CreateCard(card *models.Card, userID uint) error {
	if err := r.checkPermission(card.ListID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(card).Error; err != nil {
			return errorhandler.NewAPIError(errorhandler.ErrBadRequest, err.Error())
		}
		return nil
	})
}

func (r *GormCardRepository) GetCard(id uint, listID uint, userID uint) (*models.Card, error) {
	if err := r.checkPermission(listID, userID); err != nil {
		return nil, err
	}

	var card models.Card
	if err := r.db.Where("id = ? AND list_id = ?", id).First(&card).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NewAPIError(errorhandler.ErrNotFound, "Card not found")
		}
		return nil, errorhandler.NewGrpcInternalError()
	}

	return &card, nil
}

func (r *GormCardRepository) GetCardsByList(listID uint, userID uint) ([]*models.Card, error) {
	if err := r.checkPermission(listID, userID); err != nil {
		return nil, err
	}

	var cards []*models.Card
	if err := r.db.Where("list_id = ?", listID).Find(&cards).Error; err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	return cards, nil
}

func (r *GormCardRepository) UpdateCard(id uint, name string, listID uint, userID uint) error {
	if err := r.checkPermission(listID, userID); err != nil {
		return err
	}

	var existingCard models.Card
	if err := r.db.Where("id = ? AND list_id = ?", id, listID).First(&existingCard).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorhandler.NewAPIError(errorhandler.ErrNotFound, "Card not found")
		}
		return errorhandler.NewGrpcInternalError()
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&existingCard).Updates(existingCard).Update("name", name).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}
		return nil
	})
}

func (r *GormCardRepository) DeleteCard(id uint, listID uint, userID uint) error {
	if err := r.checkPermission(listID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND list_id = ?", id, listID).Delete(&models.Card{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(errorhandler.ErrNotFound, "Card not found")
			}
			return errorhandler.NewGrpcInternalError()
		}
		return nil
	})
}

func (r *GormCardRepository) MoveCardPosition(id uint, newPosition int, listID uint, userID uint) error {
	if err := r.checkPermission(listID, userID); err != nil {
		return err
	}

	var count int64
	r.db.Model(&models.Card{}).Where("id = ? AND list_id = ?", id, listID).Count(&count)
	if count == 0 {
		return errorhandler.NewAPIError(errorhandler.ErrNotFound, "Card not found")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get all cards
		var cards []*models.Card
		if err := tx.Where("list_id = ?", listID).Find(&cards).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(errorhandler.ErrNotFound, "Card not found")
			}
			return errorhandler.NewGrpcInternalError()
		}

		// Sort cards by position
		sort.Slice(cards, func(i, j int) bool {
			return cards[i].Position < cards[j].Position
		})

		// Update the positions of the cards
		for i, c := range cards {
			if c.ID == id {
				c.Position = newPosition
			} else {
				// Adjust the position to start from 1 and ensure no gaps
				c.Position = i + 1
				if c.Position >= newPosition {
					c.Position++
				}
			}

			if err := tx.Save(&c).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *GormCardRepository) checkPermission(listID uint, userID uint) error {
	var permission models.Permission
	if err := r.db.Where("list_id = ? AND user_id = ?", listID, userID).First(&permission).Error; err != nil {
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
