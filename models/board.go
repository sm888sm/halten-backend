package models

import (
	"gorm.io/gorm"
)

type Board struct {
	gorm.Model
	UserID      uint
	Name        string `gorm:"type:varchar(50)"`
	Visibility  string `gorm:"type:visibility_enum;default:'private'"`
	IsArchived  bool
	Permissions []Permission `gorm:"foreignKey:BoardID"`
	Lists       []List       `gorm:"foreignKey:BoardID"`
	Cards       []Card       `gorm:"foreignKey:BoardID"`
	Watches     []Watch      `gorm:"foreignKey:BoardID"`
	Labels      []Label      `gorm:"foreignKey:BoardID"`
}
