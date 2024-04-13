package models

import (
	"time"

	"gorm.io/gorm"
)

type Card struct {
	gorm.Model
	ListID      uint
	Name        string `gorm:"type:varchar(50)"`
	Description string `gorm:"type:varchar(16384)"`
	Position    int
	StartDate   *time.Time
	DueDate     *time.Time
	Attachment  Attachment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Watches     []Watch    `gorm:"foreignKey:CardID"`
}
