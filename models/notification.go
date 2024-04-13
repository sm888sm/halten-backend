package models

import (
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	ActivityLogID uint
	UserID        uint
	IsRead        bool
}
