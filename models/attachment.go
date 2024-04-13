package models

import (
	"gorm.io/gorm"
)

type Attachment struct {
	gorm.Model
	CardID    uint
	FileName  string
	FilePath  string
	Type      string `gorm:"type:type_enum;default:'document'"`
	Thumbnail string
}
