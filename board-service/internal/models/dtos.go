package model

import (
	"time"
)

type BoardDTO struct {
	ID         uint64
	Name       string
	Visibility string
	IsArchived bool
	Labels     []*LabelDTO
	Members    []*BoardMemberDTO
	Lists      []*ListMetaDTO
	Cards      []*CardMetaDTO
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type BoardMetaDTO struct {
	ID         uint64
	Name       string
	Visibility string
	IsArchived bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type LabelDTO struct {
	ID      uint64
	BoardID uint64
	Name    string
	Color   string
}

type BoardMemberDTO struct {
	ID       uint64
	Username string
	Fullname string
	Role     string
}

type ListMetaDTO struct {
	ID       uint64
	BoardID  uint
	Name     string
	Position int
}

type CardMetaDTO struct {
	ID              uint64
	ListID          uint64
	BoardID         uint64
	Name            string
	Position        int32
	Labels          []*uint64
	Members         []*uint64
	TotalAttachment uint64
	TotalComment    uint64
	IsCompleted     bool
	StartDate       *time.Time
	DueDate         *time.Time
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

type Pagination struct {
	CurrentPage  uint64
	TotalPages   uint64
	ItemsPerPage uint64
	TotalItems   uint64
	HasMore      bool
}
