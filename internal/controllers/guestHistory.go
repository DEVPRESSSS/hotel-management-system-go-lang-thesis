package controllers

import (
	"HMS-GO/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetGuestHistoryBookings(ctx *gin.Context) {

	var books []models.Book

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
		return
	}

	if err := s.Db.
		Where("books.user_id = ?", userID).
		Preload("User").
		Find(&books).Error; err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch bookings"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"reservations": books})
}
