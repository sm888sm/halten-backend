package repositories

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	models "github.com/sm888sm/halten-backend/models"

	"github.com/sm888sm/halten-backend/common/constants/roleshierarchy"
	"github.com/sm888sm/halten-backend/common/errorhandlers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) CreateUser(req *CreateUserRequest) (*CreateUserResponse, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var user models.User

		// Check if username is already in use
		result := tx.Model(&models.User{}).Where("username = ?", req.User.Username).First(&user)
		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errorhandlers.NewGrpcInternalError()
			}
		} else {
			return status.Errorf(codes.AlreadyExists, errorhandlers.NewAPIError(http.StatusConflict, "username already in use").Error())
		}

		// Check if email is already in use
		result = tx.Model(&models.User{}).Where("email = ? OR new_email = ?", req.User.Email, req.User.Email).First(&user)
		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errorhandlers.NewGrpcInternalError()
			}
		} else {
			return status.Errorf(codes.AlreadyExists, errorhandlers.NewAPIError(http.StatusConflict, "email already in use").Error())
		}

		// Create the user
		if err := tx.Create(req.User).Error; err != nil {
			return errorhandlers.NewGrpcInternalError()
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{User: req.User}, nil
}

func (r *GormUserRepository) GetUserByID(req *GetUserByIDRequest) (*GetUserByIDResponse, error) {
	var user models.User
	result := r.db.First(&user, req.UserID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errorhandlers.NewGrpcNotFoundError("User not found")
		}
		return nil, errorhandlers.NewGrpcInternalError()
	}
	return &GetUserByIDResponse{User: &user}, nil
}

func (r *GormUserRepository) GetUserByUsername(req *GetUserByUsernameRequest) (*GetUserByUsernameResponse, error) {
	var user models.User
	result := r.db.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errorhandlers.NewGrpcNotFoundError("User not found")
		}
		return nil, errorhandlers.NewGrpcInternalError()
	}
	return &GetUserByUsernameResponse{User: &user}, nil
}

func (r *GormUserRepository) UpdatePassword(req *UpdatePasswordRequest) error {
	result := r.db.Model(&models.User{}).Where("id = ?", req.UserID).Update("password", req.NewPassword)
	if result.Error != nil {
		return errorhandlers.NewGrpcInternalError()
	}
	if result.RowsAffected == 0 {
		return errorhandlers.NewGrpcNotFoundError("User not found")
	}
	return nil
}

func (r *GormUserRepository) UpdateEmail(req *UpdateEmailRequest) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var user models.User

		// Get the user
		result := tx.First(&user, req.UserID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errorhandlers.NewGrpcNotFoundError("User not found")
			}
			return errorhandlers.NewGrpcInternalError()
		}

		// Check if newEmail is already in use
		result = tx.Model(&models.User{}).Where("email = ? OR new_email = ?", req.NewEmail, req.NewEmail).First(&user)
		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errorhandlers.NewGrpcInternalError()
			}
		} else {
			return status.Errorf(codes.AlreadyExists, errorhandlers.NewAPIError(http.StatusConflict, "email already in use").Error())
		}

		// Generate a token for email verification
		token := generateToken()

		// Update email and token
		user.NewEmail = req.NewEmail
		user.Token = token
		user.TokenCreatedAt = time.Now()
		if err := tx.Save(&user).Error; err != nil {
			return errorhandlers.NewGrpcInternalError()
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *GormUserRepository) UpdateUsername(req *UpdateUsernameRequest) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var user models.User

		// Get the user
		result := tx.First(&user, req.UserID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errorhandlers.NewGrpcNotFoundError("User not found")
			}
			return errorhandlers.NewGrpcInternalError()
		}

		// Check if newUsername is already in use
		result = tx.Model(&models.User{}).Where("username = ?", req.Username).First(&user)
		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errorhandlers.NewGrpcInternalError()
			}
		} else {
			return status.Errorf(codes.AlreadyExists, errorhandlers.NewAPIError(http.StatusConflict, "Username already in use").Error())
		}

		// Update username
		user.Username = req.Username
		if err := tx.Save(&user).Error; err != nil {
			return errorhandlers.NewGrpcInternalError()
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *GormUserRepository) ConfirmEmail(req *ConfirmEmailRequest) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var user models.User

		// Get the user
		result := tx.First(&user, req.UserID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errorhandlers.NewGrpcNotFoundError("User not found")
			}
			return errorhandlers.NewGrpcInternalError()
		}

		// Check if token is valid
		if user.Token != req.Token || time.Since(user.TokenCreatedAt) > 24*time.Hour {
			return status.Errorf(codes.Unauthenticated, errorhandlers.NewAPIError(http.StatusUnauthorized, "Invalid or expired token").Error())
		}

		// Update email
		user.Email = user.NewEmail
		user.NewEmail = ""
		user.Token = ""
		user.TokenCreatedAt = time.Time{}
		if err := tx.Save(&user).Error; err != nil {
			return errorhandlers.NewGrpcInternalError()
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *GormUserRepository) CheckBoardUserRole(req *CheckBoardUserRoleRequest) error {
	var boardMember models.BoardMember
	if err := r.db.Select("role").Where("board_id = ? AND user_id = ?", req.BoardID, req.UserID).First(&boardMember).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Errorf(codes.PermissionDenied, errorhandlers.NewAPIError(http.StatusForbidden, "Board member not found").Error())
		}
		return errorhandlers.NewGrpcInternalError()
	}

	if roleshierarchy.RoleHierarchy[boardMember.Role] < roleshierarchy.RoleHierarchy[req.RequiredRole] {
		return status.Errorf(codes.PermissionDenied, errorhandlers.NewAPIError(http.StatusForbidden, "Insufficient permissions").Error())
	}

	return nil
}

func (r *GormUserRepository) CheckBoardVisibility(req *CheckBoardVisibilityRequest) error {
	var result struct {
		Visibility string
		UserID     uint64
	}

	query := `SELECT boards.visibility, board_members.user_id 
              FROM boards 
              LEFT JOIN board_members ON boards.id = board_members.board_id AND board_members.user_id = ?
              WHERE boards.id = ?`

	if err := r.db.Raw(query, req.UserID, req.BoardID).Scan(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorhandlers.NewGrpcNotFoundError("Board not found")
		}
		return errorhandlers.NewGrpcInternalError()
	}

	if result.Visibility == "private" && result.UserID != req.UserID {
		return status.Errorf(codes.PermissionDenied, errorhandlers.NewAPIError(http.StatusForbidden, "User is not a member of the specified board").Error())
	}

	return nil
}

func generateToken() string {
	b := make([]byte, 16) // generate a token of 32 characters
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
