package responsehandler

import (
	"github.com/gin-gonic/gin"
)

type Meta struct {
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

type Pagination struct {
	CurrentPage  int  `json:"currentPage"`
	TotalPages   int  `json:"totalPages"`
	ItemsPerPage int  `json:"itemsPerPage"`
	TotalItems   int  `json:"totalItems"`
	HasMore      bool `json:"hasMore"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, gin.H{
		"meta": Meta{
			Status:  status,
			Message: message,
		},
		"data": data,
	})
}

func SuccessWithPagination(c *gin.Context, status int, message string, data interface{}, pagination *Pagination) {
	c.JSON(status, gin.H{
		"meta": Meta{
			Status:     status,
			Message:    message,
			Pagination: pagination,
		},
		"data": data,
	})
}

func Error(c *gin.Context, status int, message string, errors []ErrorResponse) {
	c.JSON(status, gin.H{
		"meta": Meta{
			Status:  status,
			Message: message,
		},
		"errors": errors,
	})
}
