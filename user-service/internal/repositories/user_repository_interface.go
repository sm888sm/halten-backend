package repositories

import (
	"github.com/sm888sm/halten-backend/models"
)

type CreateUserRequest struct {
	User *models.User
}

type CreateUserResponse struct {
	User *models.User
}

type GetUserByIDRequest struct {
	UserID uint64
}

type GetUserByIDResponse struct {
	User *models.User
}

type GetUserByUsernameRequest struct {
	Username string
}

type GetUserByUsernameResponse struct {
	User *models.User
}

type UpdatePasswordRequest struct {
	UserID      uint64
	NewPassword string
}

type UpdateEmailRequest struct {
	UserID   uint64
	NewEmail string
}

type UpdateUsernameRequest struct {
	UserID   uint64
	Username string
}

type ConfirmEmailRequest struct {
	UserID uint64
	Token  string
}

type CheckBoardUserRoleRequest struct {
	UserID       uint64
	BoardID      uint64
	RequiredRole string
}

type CheckBoardVisibilityRequest struct {
	UserID  uint64
	BoardID uint64
}

type UserRepository interface {
	CreateUser(req *CreateUserRequest) (*CreateUserResponse, error)
	GetUserByID(req *GetUserByIDRequest) (*GetUserByIDResponse, error)
	GetUserByUsername(req *GetUserByUsernameRequest) (*GetUserByUsernameResponse, error)
	UpdatePassword(req *UpdatePasswordRequest) error
	UpdateEmail(req *UpdateEmailRequest) error
	UpdateUsername(req *UpdateUsernameRequest) error
	ConfirmEmail(req *ConfirmEmailRequest) error
	CheckBoardUserRole(req *CheckBoardUserRoleRequest) error
	CheckBoardVisibility(req *CheckBoardVisibilityRequest) error
	// ... Other data access methods ...
}
