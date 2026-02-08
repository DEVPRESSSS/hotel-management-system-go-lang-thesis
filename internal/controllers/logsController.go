package controllers

import (
	"HMS-GO/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get all the roles from db
func (s *Server) GetLogs(ctx *gin.Context) {

	var logs []models.HistoryLog

	if err := s.Db.Find(&logs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, logs)

}
