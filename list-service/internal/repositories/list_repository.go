package repositories

import (
	"errors"
	"sort"

	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/constants/roles"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	models "github.com/sm888sm/halten-backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormListRepository struct {
	db *gorm.DB
}

func NewListRepository(db *gorm.DB) *GormListRepository {
	return &GormListRepository{db: db}
}

func (r *GormListRepository) CreateList(list *models.List, userID uint) error {
	if err := r.checkPermission(list.BoardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(list).Error; err != nil {
			return errorhandler.NewAPIError(httpcodes.ErrBadRequest, err.Error())
		}
		return nil
	})
}

func (r *GormListRepository) GetList(id uint, boardID uint, userID uint) (*models.List, error) {
	if err := r.checkPermission(boardID, userID); err != nil {
		return nil, err
	}

	var list models.List
	if err := r.db.Where("id = ? AND board_id = ?", id).First(&list).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NewAPIError(httpcodes.ErrNotFound, "List not found")
		}
		return nil, errorhandler.NewGrpcInternalError()
	}

	return &list, nil
}

func (r *GormListRepository) GetListsByBoard(boardID uint, userID uint) ([]*models.List, error) {
	if err := r.checkPermission(boardID, userID); err != nil {
		return nil, err
	}

	var lists []*models.List
	if err := r.db.Where("board_id = ?", boardID).Find(&lists).Error; err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	return lists, nil
}

func (r *GormListRepository) UpdateList(id uint, name string, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	var existingList models.List
	if err := r.db.Where("id = ? AND board_id = ?", id, boardID).First(&existingList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorhandler.NewAPIError(httpcodes.ErrNotFound, "List not found")
		}
		return errorhandler.NewGrpcInternalError()
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&existingList).Updates(existingList).Update("name", name).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}
		return nil
	})
}

func (r *GormListRepository) DeleteList(id uint, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND board_id = ?", id, boardID).Delete(&models.List{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(httpcodes.ErrNotFound, "List not found")
			}
			return errorhandler.NewGrpcInternalError()
		}
		return nil
	})
}

func (r *GormListRepository) MoveListPosition(id uint, newPosition int, boardID uint, userID uint) error {
	if err := r.checkPermission(boardID, userID); err != nil {
		return err
	}

	var count int64
	r.db.Model(&models.List{}).Where("id = ? AND board_id = ?", id, boardID).Count(&count)
	if count == 0 {
		return errorhandler.NewAPIError(httpcodes.ErrNotFound, "List not found")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get all lists
		var lists []*models.List
		if err := tx.Where("board_id = ?", boardID).Find(&lists).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(httpcodes.ErrNotFound, "List not found")
			}
			return errorhandler.NewGrpcInternalError()
		}

		// Sort lists by position
		sort.Slice(lists, func(i, j int) bool {
			return lists[i].Position < lists[j].Position
		})

		// Update the positions of the lists
		for i, l := range lists {
			if l.ID == id {
				l.Position = newPosition
			} else {
				// Adjust the position to start from 1 and ensure no gaps
				l.Position = i + 1
				if l.Position >= newPosition {
					l.Position++
				}
			}

			if err := tx.Save(&l).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *GormListRepository) checkPermission(boardID uint, userID uint) error {
	var permission models.Permission
	if err := r.db.Where("board_id = ? AND user_id = ?", boardID, userID).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorhandler.NewAPIError(httpcodes.ErrForbidden, "Permission not found")
		}
		return err
	}

	if permission.Role == roles.OwnerRole || permission.Role == roles.AdminRole || permission.Role == roles.MemberRole {
		return nil
	}

	return errorhandler.NewAPIError(httpcodes.ErrForbidden, "User does not have permission to perform this operation")
}
