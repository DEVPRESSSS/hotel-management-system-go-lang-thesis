package controllers

import (
	"HMS-GO/internal/models"
	"HMS-GO/internal/models/dto"
	"fmt"
	"net/http"
	"time"

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

func (s *Server) GetAllEventsReservations(ctx *gin.Context) {
	var books []models.Book
	if err := s.Db.Preload("User").Find(&books).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch bookings"})
		return
	}

	events := make([]dto.Calendar, 0)

	for _, b := range books {
		events = append(events, dto.Calendar{
			Start:      b.CheckInDate.Format(time.RFC3339),
			End:        b.CheckOutDate.Format(time.RFC3339),
			Display:    "background",
			Background: "#ef4444",
		})
		events = append(events, dto.Calendar{
			Title:     b.BookId + " " + b.RoomNumber + " (" + b.User.FullName + ")",
			Start:     b.CheckInDate.Format(time.RFC3339),
			End:       b.CheckOutDate.Format(time.RFC3339),
			Color:     "#ef4444",
			TextColor: "#ffffff",
		})
	}

	ctx.JSON(http.StatusOK, events)
}

// Assign cleaner
func (s *Server) AssignCleaner(ctx *gin.Context) {

	bookingId := ctx.Param("id")

	var cleaningTask models.CleaningTask
	if err := ctx.ShouldBindJSON(&cleaningTask); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	if err := s.Db.Create(&cleaningTask).Error; err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request!!"})
		return
	}

	var reservation models.Book

	if err := ctx.ShouldBindJSON(&reservation); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}
	if err := s.Db.Model(&models.Book{}).Where("book_id = ?", bookingId).
		Update("status", "cleaning").Error; err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request!!"})
		return
	}

}

// Assign cleaner
func (s *Server) CheckinStatus(ctx *gin.Context) {

	bookingId := ctx.Param("id")

	if err := s.Db.Model(&models.Book{}).Where("book_id = ?", bookingId).
		Update("status", "check-in").Error; err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request!!"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Checking in...."})

}
