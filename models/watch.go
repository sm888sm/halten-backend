package models

import (
	"gorm.io/gorm"
)

type Watch struct {
	gorm.Model
	UserID  uint
	BoardID *uint
	ListID  *uint
	CardID  *uint
}
