package models

import (
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	BoardID uint   `gorm:"uniqueIndex:user_board_idx"`
	UserID  uint   `gorm:"uniqueIndex:user_board_idx"`
	Role    string `gorm:"type:role_enum;default:'observer'"`
}
