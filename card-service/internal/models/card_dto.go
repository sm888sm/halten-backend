package model

import (
	"time"
)

type CardDTO struct {
	ID          uint64
	ListID      uint64
	BoardID     uint64
	Name        string
	Position    int64
	Labels      []uint64
	Members     []uint64
	Attachments []uint64
	IsCompleted bool
	StartDate   *time.Time
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CardMetaDTO struct {
	ID              uint64
	ListID          uint64
	BoardID         uint64
	Name            string
	Position        int64
	Labels          []uint64
	Members         []uint64
	TotalAttachment uint64
	TotalComment    uint64
	IsCompleted     bool
	StartDate       *time.Time
	DueDate         *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
