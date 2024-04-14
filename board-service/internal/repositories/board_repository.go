package repositories

import (
	"errors"
	"math"

	internal_models "github.com/sm888sm/halten-backend/board-service/internal/models"
	"github.com/sm888sm/halten-backend/common"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	models "github.com/sm888sm/halten-backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormBoardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) *GormBoardRepository {
	return &GormBoardRepository{db: db}
}

func (r *GormBoardRepository) CreateBoard(board *models.Board, userID uint) (*models.Board, error) {
	board.UserID = userID

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(board).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		permission := models.Permission{
			UserID:  userID,
			BoardID: board.ID,
			Role:    common.OwnerRole,
		}

		if err := tx.Create(&permission).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return board, nil
}

func (r *GormBoardRepository) GetBoardByID(id uint, userID uint) (*models.Board, error) {
	var board models.Board
	var permission models.Permission

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Preload("Permissions").First(&board, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(errorhandler.ErrNotFound, "Board not found")
			}
			return errorhandler.NewGrpcInternalError()
		}

		if board.Visibility == "private" {
			if err := tx.Where("board_id = ? AND user_id = ?", id, userID).First(&permission).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errorhandler.NewAPIError(errorhandler.ErrForbidden, "You do not have permission to access this board")
				}
				return errorhandler.NewGrpcInternalError()
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &board, nil
}

func (r *GormBoardRepository) GetBoards(userID uint, pageNumber, pageSize int) (*internal_models.BoardList, error) {
	var boards []models.Board
	var totalItems int64
	offset := (pageNumber - 1) * pageSize

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN user_boards ON user_boards.board_id = boards.id").
			Where("user_boards.user_id = ?", userID).
			Offset(offset).Limit(pageSize).
			Find(&boards).Error; err != nil {
			return err
		}

		tx.Model(&models.Board{}).Joins("JOIN user_boards ON user_boards.board_id = boards.id").
			Where("user_boards.user_id = ?", userID).
			Count(&totalItems)

		return nil
	})

	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	hasMore := pageNumber < totalPages

	return &internal_models.BoardList{
		Pagination: internal_models.Pagination{
			CurrentPage:  uint64(pageNumber),
			TotalPages:   uint64(totalPages),
			ItemsPerPage: uint64(pageSize),
			TotalItems:   uint64(totalItems),
			HasMore:      hasMore,
		},
		Boards: boards,
	}, nil
}

func (r *GormBoardRepository) UpdateBoard(userID uint, id uint, name string) error {
	var board models.Board
	var permission models.Permission

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&board, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(errorhandler.ErrNotFound, "Board not found")
			}
			return errorhandler.NewGrpcInternalError()
		}

		if err := tx.Where("board_id = ? AND user_id = ?", id, userID).First(&permission).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		if permission.Role != common.AdminRole && permission.Role != common.OwnerRole {
			return errorhandler.NewAPIError(errorhandler.ErrForbidden, "Only the owner and admin can update the board")
		}

		board.Name = name
		if err := tx.Save(&board).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormBoardRepository) DeleteBoard(userID uint, id uint) error {
	var board models.Board

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&board, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(errorhandler.ErrNotFound, "Board not found")
			}
			return errorhandler.NewGrpcInternalError()
		}

		if board.UserID != userID {
			return errorhandler.NewAPIError(errorhandler.ErrForbidden, "Only the owner can delete the board")
		}

		if err := tx.Delete(&board).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormBoardRepository) GetBoardUsers(boardID uint) ([]models.User, error) {

	var permissions []models.Permission
	if err := r.db.Where("board_id = ?", boardID).Find(&permissions).Error; err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	var users []models.User
	for _, permission := range permissions {
		var user models.User
		if err := r.db.First(&user, permission.UserID).Error; err != nil {
			return nil, errorhandler.NewGrpcInternalError()
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *GormBoardRepository) AddBoardUsers(userID uint, boardID uint, userIDs []uint, role string) error {
	// Check if the current user has permission to add users
	var permission models.Permission
	if err := r.db.Where("board_id = ? AND user_id = ?", boardID, userID).First(&permission).Error; err != nil {
		return errorhandler.NewGrpcInternalError()
	}

	if permission.Role != common.AdminRole && permission.Role != common.OwnerRole {
		return errorhandler.NewAPIError(errorhandler.ErrForbidden, "Only the owner and admin can add users to the board")
	}
	existingUsersCount := 0
	for _, userID := range userIDs {
		// Check if the user already exists
		var existingPermission models.Permission
		if err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("board_id = ? AND user_id = ?", boardID, userID).First(&existingPermission).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				existingUsersCount++
				continue
			}
		} else {
			continue
		}

		permission := models.Permission{
			UserID:  userID,
			BoardID: boardID,
			Role:    role,
		}

		err := r.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Create(&permission).Error; err != nil {
				return errorhandler.NewGrpcInternalError()
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	if existingUsersCount == len(userIDs) {
		return errorhandler.NewAPIError(errorhandler.ErrConflict, "All users are already on the board")
	}

	return nil
}

func (r *GormBoardRepository) RemoveBoardUsers(userID uint, boardID uint, userIDs []uint) error {
	// Check if the current user has permission to remove users
	var permission models.Permission
	if err := r.db.Where("board_id = ? AND user_id = ?", boardID, userID).First(&permission).Error; err != nil {
		return errorhandler.NewGrpcInternalError()
	}

	if permission.Role != common.AdminRole && permission.Role != common.OwnerRole {
		return errorhandler.NewAPIError(errorhandler.ErrForbidden, "Only the owner and admin can remove users from the board")
	}

	nonExistingUsersCount := 0
	for _, userID := range userIDs {
		err := r.db.Transaction(func(tx *gorm.DB) error {
			// Check if the user is on the permission list
			var permission models.Permission
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("board_id = ? AND user_id = ?", boardID, userID).First(&permission).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					nonExistingUsersCount++
					return nil
				}
				return errorhandler.NewGrpcInternalError()
			}

			// If the user is on the permission list, proceed to delete
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Delete(&permission).Error; err != nil {
				return errorhandler.NewGrpcInternalError()
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	if nonExistingUsersCount == len(userIDs) {
		return errorhandler.NewAPIError(errorhandler.ErrNotFound, "None of the users are on the board")
	}

	return nil
}

func (r *GormBoardRepository) AssignUserRoleBoard(userID uint, boardID uint, AssignUserId uint, role string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&models.Permission{}).Where("board_id = ? AND user_id = ?", boardID, userID).Update("role", role).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})

	return err
}

func (r *GormBoardRepository) ChangeBoardOwner(boardID uint, currentOwnerID uint, newOwnerID uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var newOwnerPermission models.Permission
		if err := tx.Where("board_id = ? AND user_id = ?", boardID, newOwnerID).First(&newOwnerPermission).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandler.NewAPIError(errorhandler.ErrNotFound, "New owner must be a member of the board")
			}
			return errorhandler.NewGrpcInternalError()
		}

		var board models.Board
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&board, boardID).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		if board.UserID != currentOwnerID {
			return errorhandler.NewAPIError(errorhandler.ErrForbidden, "Only the owner can change the owner of the board")
		}

		board.UserID = newOwnerID
		if err := tx.Save(&board).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		if err := tx.Model(&models.Permission{}).Where("board_id = ? AND user_id = ?", boardID, currentOwnerID).Update("role", "member").Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		if err := tx.Model(&models.Permission{}).Where("board_id = ? AND user_id = ?", boardID, newOwnerID).Update("role", common.OwnerRole).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})

	return err
}
