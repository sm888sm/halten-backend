package model

import models "github.com/sm888sm/halten-backend/models"

type Pagination struct {
	CurrentPage  int
	TotalPages   int
	ItemsPerPage int
	TotalItems   int
	HasMore      bool
}

type BoardList struct {
	Pagination Pagination
	Boards     []models.Board
}
