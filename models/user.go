package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `gorm:"unique"`
	FullName       string
	Email          string `gorm:"unique"`
	NewEmail       string
	Password       string    `json:"-"`
	Token          string    `json:"-"`
	TokenCreatedAt time.Time `json:"-"`
	EmailConfirmed bool      `gorm:"default:false"`
	Boards         []Board   `gorm:"foreignKey:UserID"`
}
