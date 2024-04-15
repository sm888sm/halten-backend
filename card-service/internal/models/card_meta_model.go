package model

import (
	"time"

	models "github.com/sm888sm/halten-backend/models"
)

type CardMeta struct {
	ID              uint64
	ListID          uint64
	BoardID         uint64
	Name            string
	Position        int32
	Labels          []models.Label
	Members         []models.CardMember
	TotalAttachment uint64
	TotalComment    uint64
	IsCompleted     bool
	StartDate       *time.Time
	DueDate         *time.Time
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}
