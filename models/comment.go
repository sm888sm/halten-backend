package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	CardID  uint
	UserID  uint
	Content string `gorm:"type:varchar(1024)"`
	User    User   `gorm:"foreignKey:UserID"`
}
