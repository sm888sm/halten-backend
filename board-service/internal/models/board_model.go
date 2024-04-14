package model

import models "github.com/sm888sm/halten-backend/models"

type Pagination struct {
	CurrentPage  uint64
	TotalPages   uint64
	ItemsPerPage uint64
	TotalItems   uint64
	HasMore      bool
}

type BoardList struct {
	Pagination Pagination
	Boards     []models.Board
}
