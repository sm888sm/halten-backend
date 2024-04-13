package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `gorm:"unique"`
	Email          string `gorm:"unique"`
	NewEmail       string
	Password       string `json:"-"`
	Token          string `json:"-"`
	TokenCreatedAt time.Time
	EmailConfirmed bool
	Boards         []Board `gorm:"foreignKey:UserID"`
}
