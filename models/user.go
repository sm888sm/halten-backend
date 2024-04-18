package models

import (
	"time"
)

type User struct {
	BaseModel
	Username       string `gorm:"unique"`
	Fullname       string `gorm:"not null"`
	Email          string `gorm:"unique"`
	NewEmail       string
	Password       string    `json:"-"`
	Token          string    `json:"-"`
	TokenCreatedAt time.Time `json:"-"`
	EmailConfirmed bool      `gorm:"default:false"`
	Boards         []Board   `gorm:"foreignKey:UserID"`
}
