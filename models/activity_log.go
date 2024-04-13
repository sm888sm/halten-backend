package models

import (
	"gorm.io/gorm"
)

type ActivityLog struct {
	gorm.Model
	BoardID       uint
	UserID        uint
	ActionType    string `gorm:"type:varchar(50)"`
	Details       string `gorm:"type:text"`
	Notifications []Notification
}
