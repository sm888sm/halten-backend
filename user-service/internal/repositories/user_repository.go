package repositories

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	models "github.com/sm888sm/halten-backend/models"

	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/constants/roleshierarchy"
	"github.com/sm888sm/halten-backend/common/errorhandler"

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
				errorhandler.NewGrpcInternalError()
			}
		} else {
			return status.Errorf(codes.AlreadyExists, errorhandler.NewAPIError(httpcodes.ErrConflict, "username already in use").Error())
		}

		// Check if email is already in use
		result = tx.Model(&models.User{}).Where("email = ? OR new_email = ?", req.User.Email, req.User.Email).First(&user)
		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errorhandler.NewGrpcInternalError()
			}
		} else {
			return status.Errorf(codes.AlreadyExists, errorhandler.NewAPIError(httpcodes.ErrConflict, "email already in use").Error())
		}

		// Create the user
		if err := tx.Create(req.User).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
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
			return nil, status.Errorf(codes.NotFound, errorhandler.NewAPIError(httpcodes.ErrNotFound, "User not found").Error())
		}
		return nil, errorhandler.NewGrpcInternalError()
	}
	return &GetUserByIDResponse{User: &user}, nil
}

func (r *GormUserRepository) GetUserByUsername(req *GetUserByUsernameRequest) (*GetUserByUsernameResponse, error) {
	var user models.User
	result := r.db.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, errorhandler.NewAPIError(httpcodes.ErrNotFound, "User not found").Error())
		}
		return nil, errorhandler.NewGrpcInternalError()
	}
	return &GetUserByUsernameResponse{User: &user}, nil
}

func (r *GormUserRepository) UpdatePassword(req *UpdatePasswordRequest) error {
	result := r.db.Model(&models.User{}).Where("id = ?", req.UserID).Update("password", req.NewPassword)
	if result.Error != nil {
		return errorhandler.NewGrpcInternalError()
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, errorhandler.NewAPIError(httpcodes.ErrNotFound, "User not found").Error())
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
				return status.Errorf(codes.NotFound, errorhandler.NewAPIError(httpcodes.ErrNotFound, "User not found").Error())
			}
			return errorhandler.NewGrpcInternalError()
		}

		// Check if newEmail is already in use
		result = tx.Model(&models.User{}).Where("email = ? OR new_email = ?", req.NewEmail, req.NewEmail).First(&user)
		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errorhandler.NewGrpcInternalError()
			}
		} else {
			return status.Errorf(codes.AlreadyExists, errorhandler.NewAPIError(httpcodes.ErrConflict, "email already in use").Error())
		}

		// Generate a token for email verification
		token := generateToken()

		// Update email and token
		user.NewEmail = req.NewEmail
		user.Token = token
		user.TokenCreatedAt = time.Now()
		if err := tx.Save(&user).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
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
				return status.Errorf(codes.NotFound, errorhandler.NewAPIError(httpcodes.ErrNotFound, "User not found").Error())
			}
			return errorhandler.NewGrpcInternalError()
		}

		// Check if newUsername is already in use
		result = tx.Model(&models.User{}).Where("username = ?", req.Username).First(&user)
		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errorhandler.NewGrpcInternalError()
			}
		} else {
			return status.Errorf(codes.AlreadyExists, errorhandler.NewAPIError(httpcodes.ErrConflict, "Username already in use").Error())
		}

		// Update username
		user.Username = req.Username
		if err := tx.Save(&user).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *GormUserRepository) ConfirmNewEmail(req *ConfirmNewEmailRequest) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var user models.User

		// Get the user
		result := tx.First(&user, req.UserID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, errorhandler.NewAPIError(httpcodes.ErrNotFound, "User not found").Error())
			}
			return errorhandler.NewGrpcInternalError()
		}

		// Check if token is valid
		if user.Token != req.Token || time.Since(user.TokenCreatedAt) > 24*time.Hour {
			return status.Errorf(codes.InvalidArgument, errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid or expired token").Error())
		}

		// Update email
		user.Email = user.NewEmail
		user.NewEmail = ""
		user.Token = ""
		user.TokenCreatedAt = time.Time{}
		if err := tx.Save(&user).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *GormUserRepository) CheckBoardUserRole(req *CheckBoardUserRoleRequest) error {
	var permission models.Permission
	if err := r.db.Select("role").Where("board_id = ? AND user_id = ?", req.BoardID, req.UserID).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorhandler.NewAPIError(httpcodes.ErrForbidden, "Permission not found")
		}
		return errorhandler.NewGrpcInternalError()
	}

	if roleshierarchy.RoleHierarchy[permission.Role] < roleshierarchy.RoleHierarchy[req.RequiredRole] {
		return errorhandler.NewAPIError(httpcodes.ErrForbidden, "Insufficient permissions")
	}

	return nil
}

func (r *GormUserRepository) CheckBoardVisibility(req *CheckBoardVisibilityRequest) error {
	var result struct {
		Visibility string
		UserID     uint
	}

	query := `SELECT boards.visibility, permissions.user_id 
              FROM boards 
              LEFT JOIN permissions ON boards.id = permissions.board_id AND permissions.user_id = ?
              WHERE boards.id = ?`

	if err := r.db.Raw(query, req.UserID, req.BoardID).Scan(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorhandler.NewAPIError(httpcodes.ErrNotFound, "Board not found")
		}
		return errorhandler.NewGrpcInternalError()
	}

	if result.Visibility == "private" && result.UserID != uint(req.UserID) {
		return errorhandler.NewAPIError(httpcodes.ErrForbidden, "Permission not found")
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
