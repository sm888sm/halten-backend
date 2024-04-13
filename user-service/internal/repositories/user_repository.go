package repositories

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os/user"
	"time"

	models "github.com/sm888sm/halten-backend/models"

	"github.com/sm888sm/halten-backend/common/errorhandler"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) CreateUser(user *models.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if username is already in use
		var count int64
		tx.Model(user).Where("username = ?", user.Username).Count(&count)
		if count > 0 {
			return status.Errorf(codes.AlreadyExists, errorhandler.NewAPIError(errorhandler.ErrConflict, "username already in use").Error())
		}

		// Check if email is already in use
		tx.Model(user).Where("email = ? OR new_email = ?", user.Email, user.Email).Count(&count)
		if count > 0 {
			return status.Errorf(codes.AlreadyExists, errorhandler.NewAPIError(errorhandler.ErrConflict, "email already in use").Error())
		}

		// Create user
		if err := tx.Create(user).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func (r *GormUserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, errorhandler.NewAPIError(errorhandler.ErrNotFound, "User not found").Error())
		}
		return nil, errorhandler.NewGrpcInternalError()
	}
	return &user, nil
}

func (r *GormUserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, errorhandler.NewAPIError(errorhandler.ErrNotFound, "User not found").Error())
		}
		return nil, errorhandler.NewGrpcInternalError()
	}
	return &user, nil
}

func (r *GormUserRepository) UpdatePassword(userID uint, newPassword string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).First(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, errorhandler.NewAPIError(errorhandler.ErrNotFound, "User not found").Error())
			}
			return errorhandler.NewGrpcInternalError()
		}

		// Update password (assuming proper password hashing)
		user.Password = newPassword
		return tx.Save(&user).Error
	})
}

func (r *GormUserRepository) UpdateEmail(userID uint, newEmail string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).First(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, errorhandler.NewAPIError(errorhandler.ErrNotFound, "User not found").Error())
			}
			return errorhandler.NewGrpcInternalError()
		}

		// Check if newEmail is already in use
		var count int64
		tx.Model(&models.User{}).Where("email = ? OR new_email = ?", newEmail, newEmail).Count(&count)
		if count > 0 {
			return status.Errorf(codes.AlreadyExists, errorhandler.NewAPIError(errorhandler.ErrConflict, "email already in use").Error())
		}

		// Generate a token for email verification
		token := generateToken()

		// Update email and token
		user.NewEmail = newEmail
		user.Token = token
		user.TokenCreatedAt = time.Now()
		if err := tx.Save(&user).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		// TODO: Send the token to the new email address
		// sendEmail(newEmail, token)

		return nil
	})
}

func (r *GormUserRepository) UpdateUsername(oldUsername, newUsername string) error {

	return r.db.Transaction(func(tx *gorm.DB) error {

		var user user.User
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("username = ?", oldUsername).First(&user)
		if result.Error != nil {
			return errorhandler.NewGrpcInternalError()
		}

		// Check if the new username is already taken
		if _, err := r.GetUserByUsername(newUsername); err == nil {
			return status.Errorf(codes.AlreadyExists, errorhandler.NewAPIError(errorhandler.ErrNotFound, "username already taken").Error())
		}

		// Update username
		user.Username = newUsername
		return tx.Save(user).Error

	})
}

func (r *GormUserRepository) ConfirmNewEmail(username string, token string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("username = ?", username).First(&user)
		if result.Error != nil {
			return errorhandler.NewAPIError(errorhandler.ErrInternalServerError, "Internal server error")
		}

		// Check if the token is valid
		if user.Token != token {
			return status.Errorf(codes.Unauthenticated, errorhandler.NewAPIError(errorhandler.ErrUnauthorized, "Invalid token").Error())
		}

		// Check if the token is expired
		if time.Since(user.TokenCreatedAt) > 48*time.Hour {
			return status.Errorf(codes.Unauthenticated, errorhandler.NewAPIError(errorhandler.ErrUnauthorized, "Token expired").Error())
		}

		// Update email
		user.Email = user.NewEmail
		user.NewEmail = ""
		user.Token = ""
		user.TokenCreatedAt = time.Time{} // reset the token creation time
		if err := tx.Save(&user).Error; err != nil {
			return errorhandler.NewGrpcInternalError()
		}

		return nil
	})
}

func generateToken() string {
	b := make([]byte, 16) // generate a token of 32 characters
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
