package models

import (
	"gorm.io/gorm"
)

type List struct {
	gorm.Model
	BoardID  uint   `gorm:"foreign_key"`
	Name     string `gorm:"type:varchar(50)"`
	Position int
	Cards    []Card  `gorm:"foreignKey:ListID"`
	Watches  []Watch `gorm:"foreignKey:ListID"`
}
