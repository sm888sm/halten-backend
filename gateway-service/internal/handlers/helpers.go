package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sm888sm/halten-backend/common/errorhandlers"
)

func getUserID(c *gin.Context) (uint64, error) {
	user, exists := c.Get("userID")
	if !exists {
		return 0, errorhandlers.NewHttpInternalError()
	}

	userID, ok := user.(uint64)
	if !ok {
		return 0, errorhandlers.NewHttpInternalError()
	}

	return userID, nil
}
