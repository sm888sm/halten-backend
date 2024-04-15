package models

import (
	"gorm.io/gorm"
)

type Attachment struct {
	gorm.Model
	BoardID   uint `gorm:"foreignKey:CardID"`
	CardID    uint `gorm:"foreignKey:CardID"`
	FileName  string
	FilePath  string
	Type      string `gorm:"type:type_enum;default:'document'"`
	Thumbnail string
}
