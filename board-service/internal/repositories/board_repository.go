package repositories

import (
	"errors"

	dtos "github.com/sm888sm/halten-backend/board-service/internal/models"
	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/constants/roles"
	"github.com/sm888sm/halten-backend/common/errorhandlers"
	"github.com/sm888sm/halten-backend/common/helpers"
	"github.com/sm888sm/halten-backend/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gorm.io/gorm"
)

type GormBoardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) *GormBoardRepository {
	return &GormBoardRepository{db: db}
}

func (r *GormBoardRepository) CreateBoard(req *CreateBoardRequest) (*CreateBoardResponse, error) {
	board := *req.Board
	err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&board).Error; err != nil {
			return errorhandlers.NewGrpcInternalError()
		}

		boardMember := models.BoardMember{
			UserID:  req.Board.UserID,
			BoardID: board.ID,
			Role:    roles.OwnerRole,
		}

		if err := tx.Create(&boardMember).Error; err != nil {
			return errorhandlers.NewGrpcInternalError()
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &CreateBoardResponse{Board: &board}, nil
}

func (r *GormBoardRepository) GetBoardByID(req *GetBoardByIDRequest) (*GetBoardByIDResponse, error) {
	var board models.Board
	var boardDTO dtos.BoardDTO

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Select("ID", "Name", "Visibility", "IsArchived").
			Preload("Labels", func(db *gorm.DB) *gorm.DB {
				return db.Select("ID", "Name", "Color")
			}).
			Preload("Members", func(db *gorm.DB) *gorm.DB {
				return db.Select("ID", "Role")
			}).
			First(&board, req.BoardID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorhandlers.NewAPIError(httpcodes.ErrNotFound, "Board not found")
			}
			return errorhandlers.NewGrpcInternalError()
		}

		// Map the board to the BoardDTO struct
		boardDTO.ID = board.ID
		boardDTO.Name = board.Name
		boardDTO.Visibility = board.Visibility
		boardDTO.IsArchived = board.IsArchived

		// Map the labels
		for _, label := range board.Labels {
			boardDTO.Labels = append(boardDTO.Labels, &dtos.LabelDTO{
				ID:    label.ID,
				Name:  label.Name,
				Color: label.Color,
			})
		}

		// Map the members
		for _, member := range board.Members {
			var user models.User
			if err := tx.First(&user, member.UserID).Error; err != nil {
				return errorhandlers.NewGrpcInternalError()
			}

			boardDTO.Members = append(boardDTO.Members, &dtos.BoardMemberDTO{
				ID:       user.ID,
				Username: user.Username,
				Fullname: user.Fullname,
				Role:     member.Role,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &GetBoardByIDResponse{Board: boardDTO}, nil
}

func (r *GormBoardRepository) GetBoardList(req *GetBoardListRequest) (*GetBoardListResponse, error) {
	pageNumber, pageSize := int(req.PageNumber), int(req.PageSize)
	var boards []*models.Board
	var boardMetaDTOs []*dtos.BoardMetaDTO
	var totalItems int64

	offset := (pageNumber - 1) * pageSize

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("boards.id, boards.created_at, boards.updated_at, boards.name, boards.visibility").
			Joins("JOIN board_members ON board_members.board_id = boards.id").
			Where("board_members.user_id = ? AND boards.is_archived = ?", req.UserID, false).
			Offset(offset).Limit(pageSize).
			Find(&boards).Error; err != nil {
			return err
		}

		// Map the boards to the BoardMetaDTO struct
		for _, board := range boards {
			boardMetaDTOs = append(boardMetaDTOs, &dtos.BoardMetaDTO{
				ID:         board.ID,
				Name:       board.Name,
				Visibility: board.Visibility,
				IsArchived: board.IsArchived,
			})
		}

		// Count the total items
		tx.Model(&models.Board{}).
			Joins("JOIN board_members ON board_members.board_id = boards.id").
			Where("board_members.user_id = ? AND boards.is_archived = ?", req.UserID, false).
			Count(&totalItems)

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Calculate the pagination data
	pagination := &dtos.Pagination{
		CurrentPage:  req.PageNumber,
		TotalPages:   (uint64(totalItems) + req.PageSize - 1) / req.PageSize,
		ItemsPerPage: req.PageSize,
		TotalItems:   uint64(totalItems),
		HasMore:      req.PageNumber*req.PageSize < uint64(totalItems),
	}

	return &GetBoardListResponse{
		Boards:     boardMetaDTOs,
		Pagination: pagination,
	}, nil
}

func (r *GormBoardRepository) UpdateBoardName(req *UpdateBoardNameRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Board{}).Where("id = ?", req.BoardID).Update("name", req.Name)
		if result.Error != nil {
			return errorhandlers.NewGrpcInternalError()
		}
		if result.RowsAffected == 0 {
			return status.Errorf(codes.NotFound, errorhandlers.NewAPIError(httpcodes.ErrNotFound, "Board not found").Error())
		}
		return nil
	})
}

func (r *GormBoardRepository) GetBoardMembers(req *GetBoardMembersRequest) (*GetBoardMembersResponse, error) {
	var boardMemberDTOs []*dtos.BoardMemberDTO

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var boardMembers []models.BoardMember

		if err := tx.Where("board_id = ?", req.BoardID).Find(&boardMembers).Error; err != nil {
			return errorhandlers.NewGrpcInternalError()
		}

		for _, boardMember := range boardMembers {
			var user models.User
			if err := tx.First(&user, boardMember.UserID).Error; err != nil {
				return errorhandlers.NewGrpcInternalError()
			}

			boardMemberDTOs = append(boardMemberDTOs, &dtos.BoardMemberDTO{
				ID:       user.ID,
				Username: user.Username,
				Role:     boardMember.Role,
				Fullname: user.Fullname,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &GetBoardMembersResponse{
		Members: boardMemberDTOs,
	}, nil
}

func (r *GormBoardRepository) AddBoardUsers(req *AddBoardUsersRequest) error {
	// Validate all users before starting the transaction
	var existingUserIDs []uint
	if err := r.db.Model(&models.User{}).Where("id IN ?", req.UserIDs).Select("id").Find(&existingUserIDs).Error; err != nil {
		return errorhandlers.NewGrpcInternalError()
	}

	if len(existingUserIDs) != len(req.UserIDs) {
		return status.Errorf(codes.NotFound, errorhandlers.NewAPIError(httpcodes.ErrNotFound, "One or more users not found").Error())
	}

	// Check if any user is already a member of the board
	var existingBoardMemberIDs []uint
	if err := r.db.Model(&models.BoardMember{}).Where("board_id = ? AND user_id IN ?", req.BoardID, req.UserIDs).Select("id").Find(&existingBoardMemberIDs).Error; err != nil {
		return errorhandlers.NewGrpcInternalError()
	}

	if len(existingBoardMemberIDs) > 0 {
		return errorhandlers.NewAPIError(httpcodes.ErrConflict, "One or more users are already members of the board")
	}

	// Start the transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, userID := range req.UserIDs {
			boardMember := models.BoardMember{
				UserID:  userID,
				BoardID: req.BoardID,
				Role:    req.Role,
			}

			// Add the user to the board
			if err := tx.Create(&boardMember).Error; err != nil {
				return errorhandlers.NewGrpcInternalError()
			}
		}

		return nil
	})
}

func (r *GormBoardRepository) RemoveBoardUsers(req *RemoveBoardUsersRequest) error {
	// Check if all users exist
	var existingUserIDs []uint
	if err := r.db.Model(&models.User{}).Where("id IN ?", req.UserIDs).Select("id").Find(&existingUserIDs).Error; err != nil {
		return errorhandlers.NewGrpcInternalError()
	}

	if len(existingUserIDs) != len(req.UserIDs) {
		return status.Errorf(codes.NotFound, errorhandlers.NewAPIError(httpcodes.ErrNotFound, "One or more users not found").Error())
	}

	// Check if all users are members of the board
	var existingBoardMemberIDs []uint
	if err := r.db.Model(&models.BoardMember{}).Where("board_id = ? AND user_id IN ?", req.BoardID, req.UserIDs).Select("id").Find(&existingBoardMemberIDs).Error; err != nil {
		return errorhandlers.NewGrpcInternalError()
	}

	if len(existingBoardMemberIDs) != len(req.UserIDs) {
		return status.Errorf(codes.AlreadyExists, errorhandlers.NewAPIError(httpcodes.ErrConflict, "One or more users are not members of the board").Error())
	}

	// Start the transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("board_id = ? AND user_id IN ?", req.BoardID, req.UserIDs).Delete(&models.BoardMember{}).Error; err != nil {
			return errorhandlers.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormBoardRepository) AssignBoardUsersRole(req *AssignBoardUsersRoleRequest) error {
	// Check if all users exist and are members of the board
	var existingBoardMemberIDs []uint
	if err := r.db.Model(&models.BoardMember{}).Where("board_id = ? AND user_id IN ?", req.BoardID, req.UserIDs).Select("id").Find(&existingBoardMemberIDs).Error; err != nil {
		return errorhandlers.NewGrpcInternalError()
	}

	if len(existingBoardMemberIDs) != len(req.UserIDs) {
		return status.Errorf(codes.NotFound, errorhandlers.NewAPIError(httpcodes.ErrNotFound, "One or more users not found or not members of the board").Error())
	}

	// Check if the current user has permission to assign roles
	var currentUserRole models.BoardMember
	if err := r.db.Where("board_id = ? AND user_id = ?", req.BoardID, req.UserID).Select("role").First(&currentUserRole).Error; err != nil {
		return errorhandlers.NewGrpcInternalError()
	}

	// Check if the current user can assign the requested role
	if !canAssignRole(currentUserRole.Role, req.Role) {
		return status.Errorf(codes.PermissionDenied, errorhandlers.NewAPIError(httpcodes.ErrForbidden, "You don't have permission to assign this role").Error())
	}

	// Start the transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.BoardMember{}).Where("board_id = ? AND user_id IN ?", req.BoardID, req.UserIDs).Update("role", req.Role).Error; err != nil {
			return errorhandlers.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormBoardRepository) ChangeBoardOwner(req *ChangeBoardOwnerRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if the new owner exists
		var newUser models.User
		if err := tx.First(&newUser, req.NewOwnerID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, errorhandlers.NewAPIError(httpcodes.ErrNotFound, "User not found").Error())
			}
			return errorhandlers.NewGrpcInternalError()
		}

		// Check if the new owner is already a member of the board
		var newOwnerMember models.BoardMember
		if err := tx.Where("board_id = ? AND user_id = ?", req.BoardID, req.NewOwnerID).First(&newOwnerMember).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, errorhandlers.NewAPIError(httpcodes.ErrForbidden, "User is not a member of the board").Error())
			}
			return errorhandlers.NewGrpcInternalError()
		}

		// Check if the current user is the owner of the board
		var currentOwnerMember models.BoardMember
		if err := tx.Where("board_id = ? AND user_id = ? AND role = ?", req.BoardID, req.CurrentOwnerID, roles.OwnerRole).First(&currentOwnerMember).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.PermissionDenied, errorhandlers.NewAPIError(httpcodes.ErrForbidden, "Current user is not the owner of the board").Error())
			}
			return errorhandlers.NewGrpcInternalError()
		}

		// Update the role of the current owner to 'member'
		if err := tx.Model(&currentOwnerMember).Update("role", roles.MemberRole).Error; err != nil {
			return errorhandlers.NewGrpcInternalError()
		}

		// Update the role of the new owner to 'owner'
		if err := tx.Model(&newOwnerMember).Update("role", roles.OwnerRole).Error; err != nil {
			return errorhandlers.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormBoardRepository) ChangeBoardVisibility(req *ChangeBoardVisibilityRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if the board exists
		var board models.Board
		if err := tx.First(&board, req.BoardID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, errorhandlers.NewAPIError(httpcodes.ErrNotFound, "Board not found").Error())
			}
			return errorhandlers.NewGrpcInternalError()
		}

		// Validate the visibility value
		validVisibilities := []string{"private", "public", "team"}
		if !helpers.Contains(validVisibilities, req.Visibility) {
			return status.Errorf(codes.InvalidArgument, errorhandlers.NewAPIError(httpcodes.ErrBadRequest, "Invalid visibility value").Error())
		}

		// Update the visibility of the board
		result := tx.Model(&models.Board{}).Where("id = ?", req.BoardID).Update("visibility", req.Visibility)
		if result.Error != nil {
			return errorhandlers.NewGrpcInternalError()
		}
		if result.RowsAffected == 0 {
			return errorhandlers.NewAPIError(httpcodes.ErrNotFound, "Board not found")
		}
		return nil
	})
}

func (r *GormBoardRepository) AddLabel(req *AddLabelRequest) (*models.Label, error) {
	label := models.Label{
		Name:  req.Name,
		Color: req.Color,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&label).Error; err != nil {
			return errorhandlers.NewGrpcInternalError()
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &label, nil
}

func (r *GormBoardRepository) RemoveLabel(req *RemoveLabelRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Find the label
		var label models.Label
		if err := tx.Where("id = ? AND board_id = ?", req.LabelID, req.BoardID).First(&label).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, errorhandlers.NewAPIError(httpcodes.ErrNotFound, "Label not found").Error())
			}
			return errorhandlers.NewGrpcInternalError()
		}

		// Remove the label from all cards
		if err := tx.Model(&label).Association("Cards").Clear(); err != nil {
			return errorhandlers.NewGrpcInternalError()
		}

		// Delete the label
		if err := tx.Delete(&models.Label{}, req.LabelID).Error; err != nil {
			return errorhandlers.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormBoardRepository) GetArchivedBoardList(req *GetArchivedBoardListRequest) (*GetArchivedBoardListResponse, error) {
	pageNumber, pageSize := int(req.PageNumber), int(req.PageSize)
	var boards []*models.Board
	var boardMetaDTOs []*dtos.BoardMetaDTO
	var totalItems int64

	offset := (pageNumber - 1) * pageSize

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("boards.id, boards.created_at, boards.updated_at, boards.name, boards.visibility").
			Joins("JOIN board_members ON board_members.board_id = boards.id").
			Where("board_members.user_id = ? AND boards.is_archived = ?", req.UserID, true).
			Offset(offset).Limit(pageSize).
			Find(&boards).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, errorhandlers.NewAPIError(httpcodes.ErrNotFound, "Board not found").Error())
			}
			return errorhandlers.NewGrpcInternalError()
		}

		// Map the boards to the BoardMetaDTO struct
		for _, board := range boards {
			boardMetaDTOs = append(boardMetaDTOs, &dtos.BoardMetaDTO{
				ID:         board.ID,
				Name:       board.Name,
				Visibility: board.Visibility,
				IsArchived: board.IsArchived,
			})
		}

		// Count the total items
		if err := tx.Model(&models.Board{}).
			Joins("JOIN board_members ON board_members.board_id = boards.id").
			Where("board_members.user_id = ? AND boards.is_archived = ?", req.UserID, false).
			Count(&totalItems).Error; err != nil {
			return errorhandlers.NewGrpcInternalError()
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Calculate the pagination data
	pagination := dtos.Pagination{
		CurrentPage:  req.PageNumber,
		TotalPages:   (uint64(totalItems) + req.PageSize - 1) / req.PageSize,
		ItemsPerPage: req.PageSize,
		TotalItems:   uint64(totalItems),
		HasMore:      req.PageNumber*req.PageSize < uint64(totalItems),
	}

	return &GetArchivedBoardListResponse{
		Boards:     boardMetaDTOs,
		Pagination: &pagination,
	}, nil
}

func (r *GormBoardRepository) ArchiveBoard(req *ArchiveBoardRequest) error {
	// Update the 'archived' field of the board to true
	if err := r.db.Model(&models.Board{}).Where("id = ?", req.BoardID).Update("archived", true).Error; err != nil {
		return errorhandlers.NewGrpcInternalError()
	}

	return nil
}

func (r *GormBoardRepository) RestoreBoard(req *RestoreBoardRequest) error {
	// Update the 'archived' field of the board to false
	if err := r.db.Model(&models.Board{}).Where("id = ?", req.BoardID).Update("archived", false).Error; err != nil {
		return errorhandlers.NewGrpcInternalError()
	}

	return nil
}

func (r *GormBoardRepository) DeleteBoard(req *DeleteBoardRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ? AND is_archived = ?", req.BoardID, true).Delete(&models.Board{})
		if result.Error != nil {
			return errorhandlers.NewGrpcInternalError()
		}
		if result.RowsAffected == 0 {
			return status.Errorf(codes.NotFound, errorhandlers.NewAPIError(httpcodes.ErrNotFound, "Board not found or not archived").Error())
		}
		return nil
	})
}

func (r *GormBoardRepository) GetBoardIDByList(req *GetBoardIDByListRequest) (uint64, error) {
	var list models.List
	err := r.db.First(&list, req.ListID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errorhandlers.NewAPIError(httpcodes.ErrNotFound, "List not found")
		}
		return 0, errorhandlers.NewGrpcInternalError()
	}
	return list.BoardID, nil
}

func (r *GormBoardRepository) GetBoardIDByCard(req *GetBoardIDByCardRequest) (uint64, error) {
	var card models.Card
	err := r.db.First(&card, req.CardID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errorhandlers.NewAPIError(httpcodes.ErrNotFound, "Card not found")
		}
		return 0, errorhandlers.NewGrpcInternalError()
	}
	return card.BoardID, nil
}
