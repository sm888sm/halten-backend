package models

import (
	"time"

	"gorm.io/gorm"
)

type Card struct {
	gorm.Model
	BoardID     uint   `gorm:"not null"`
	ListID      uint   `gorm:"not null"`
	Name        string `gorm:"type:varchar(50);not null"`
	Description string `gorm:"type:varchar(16384)"`
	Position    int
	IsArchived  bool
	IsCompleted bool
	StartDate   *time.Time
	DueDate     *time.Time
	Attachments []Attachment `gorm:"foreignKey:CardID"`
	Comments    []Comment    `gorm:"foreignKey:CardID"`
	Labels      []Label      `gorm:"foreignKey:CardID"`
	Members     []CardMember `gorm:"foreignKey:CardID"`
	Watches     []Watch      `gorm:"foreignKey:CardID"`
}
