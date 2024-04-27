package repositories

import (
	"errors"
	"sort"

	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
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

func (r *GormListRepository) CreateList(req *CreateListRequest) (*CreateListResponse, error) {

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(req.List).Error; err != nil {
			return errorhandler.NewAPIError(httpcodes.ErrBadRequest, err.Error())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &CreateListResponse{List: req.List}, nil
}

func (r *GormListRepository) GetListByID(req *GetListRequest) (*GetListResponse, error) {

	var list models.List
	if err := r.db.Where("id = ? AND board_id = ?", req.ID, req.BoardID).First(&list).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NewAPIError(httpcodes.ErrNotFound, "List not found")
		}
		return nil, errorhandler.NewGrpcInternalError()
	}

	return &GetListResponse{List: &list}, nil
}

func (r *GormListRepository) GetListsByBoard(req *GetListsByBoardRequest) (*GetListsByBoardResponse, error) {

	var lists []*models.List
	if err := r.db.Where("board_id = ?", req.BoardID).Find(&lists).Error; err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	return &GetListsByBoardResponse{Lists: lists}, nil
}

func (r *GormListRepository) UpdateListName(req *UpdateListNameRequest) error {

	var existingList models.List
	if err := r.db.Where("id = ? AND board_id = ?", req.ID, req.BoardID).First(&existingList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorhandler.NewAPIError(httpcodes.ErrNotFound, "List not found")
		}
		return errorhandler.NewGrpcInternalError()
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&existingList).Updates(existingList).Update("name", req.Name).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}
		return nil
	})
}

func (r *GormListRepository) MoveListPosition(req *MoveListPositionRequest) error {

	var count int64
	r.db.Model(&models.List{}).Where("id = ? AND board_id = ?", req.ID, req.BoardID).Count(&count)
	if count == 0 {
		return errorhandler.NewAPIError(httpcodes.ErrNotFound, "List not found")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get all lists
		var lists []*models.List
		if err := tx.Where("board_id = ?", req.BoardID).Find(&lists).Error; err != nil {
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
			if l.ID == req.ID {
				l.Position = req.NewPosition
			} else {
				// Adjust the position to start from 1 and ensure no gaps
				l.Position = int64(i + 1)
				if l.Position >= req.NewPosition {
					l.Position++
				}
			}
		}

		// Save the updated lists
		for _, l := range lists {
			if err := tx.Save(l).Error; err != nil {
				return errorhandler.NewGrpcInternalError()
			}
		}

		return nil
	})
}

func (r *GormListRepository) ArchiveList(req *ArchiveListRequest) error {
	err := r.db.Model(&models.List{}).Where("id = ? AND board_id = ?", req.ListID, req.BoardID).Update("archived", true).Error
	if err != nil {
		return errorhandler.NewGrpcInternalError()
	}
	return nil
}

func (r *GormListRepository) RestoreList(req *RestoreListRequest) error {
	err := r.db.Model(&models.List{}).Where("id = ? AND board_id = ?", req.ListID, req.BoardID).Update("archived", false).Error
	if err != nil {
		return errorhandler.NewGrpcInternalError()
	}
	return nil
}

func (r *GormListRepository) DeleteList(req *DeleteListRequest) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND board_id = ?", req.ID, req.BoardID).Delete(&models.List{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(httpcodes.ErrNotFound, "List not found")
			}
			return errorhandler.NewGrpcInternalError()
		}
		return nil
	})
}
