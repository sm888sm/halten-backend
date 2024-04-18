package models

import (
	"time"
)

type Card struct {
	BaseModel
	BoardID     uint64 `gorm:"not null"`
	ListID      uint64 `gorm:"not null"`
	Name        string `gorm:"type:varchar(50);not null"`
	Description string `gorm:"type:varchar(16384)"`
	Position    int64  `gorm:"not null"`
	IsArchived  bool
	IsCompleted bool
	StartDate   *time.Time
	DueDate     *time.Time
	Attachments []Attachment `gorm:"foreignKey:CardID"`
	Comments    []Comment    `gorm:"foreignKey:CardID"`
	Labels      []Label      `gorm:"many2many:card_labels"`
	Members     []CardMember `gorm:"foreignKey:CardID"`
	Watches     []Watch      `gorm:"foreignKey:CardID"`
}
