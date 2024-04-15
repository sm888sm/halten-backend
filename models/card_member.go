package models

import (
	"gorm.io/gorm"
)

type CardMember struct {
	gorm.Model
	CardID uint
	UserID uint
	// other fields as needed
}
