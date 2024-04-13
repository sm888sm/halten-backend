package repositories

import (
	"github.com/sm888sm/halten-backend/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdatePassword(userID uint, newPassword string) error
	UpdateEmail(userID uint, newEmail string) error
	UpdateUsername(oldUsername, newUsername string) error
	ConfirmNewEmail(username string, token string) error
	// ... Other data access methods ...
}
