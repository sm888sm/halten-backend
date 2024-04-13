package models

import (
	"gorm.io/gorm"
)

type Board struct {
	gorm.Model
	UserID      uint
	Name        string       `gorm:"type:varchar(50)"`
	Visibility  string       `gorm:"type:visibility_enum;default:'private'"`
	Permissions []Permission `gorm:"foreignKey:BoardID"`
	Lists       []List       `gorm:"foreignKey:BoardID"`
	Watches     []Watch      `gorm:"foreignKey:BoardID"`
}
