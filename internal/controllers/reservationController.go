package controllers

import (
	"HMS-GO/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetAllReservations(ctx *gin.Context) {

	var books []models.Book
	if err := s.Db.Preload("User").
		Find(&books).Error; err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch bookings"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reservations": books})

}
