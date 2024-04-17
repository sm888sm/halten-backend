package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sm888sm/halten-backend/common/errorhandler"
)

func getUserID(c *gin.Context) (uint64, error) {
	user, exists := c.Get("userID")
	if !exists {
		return 0, errorhandler.NewHttpInternalError()
	}

	userID, ok := user.(uint64)
	if !ok {
		return 0, errorhandler.NewHttpInternalError()
	}

	return userID, nil
}
