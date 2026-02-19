package controllers

import (
	"HMS-GO/internal/models"
	"HMS-GO/internal/models/dto"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	var payload struct {
		CleanerIds []string `json:"cleaner_id"`
		RoomId     string   `json:"room_id"`
	}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload", "details": err.Error()})
		return
	}

	// Validate that at least one cleaner is selected
	if len(payload.CleanerIds) == 0 {
		ctx.JSON(400, gin.H{"error": "At least one cleaner must be selected"})
		return
	}

	// Create cleaning tasks for each cleaner
	for _, cleanerId := range payload.CleanerIds {

		cleaningTask := models.CleaningTask{
			Id:        uuid.New().String(),
			BookId:    bookingId,
			RoomId:    payload.RoomId,
			CleanerId: &cleanerId,
			Status:    "Cleaning in progress",
		}

		if err := s.Db.Create(&cleaningTask).Error; err != nil {
			fmt.Printf("Error creating cleaning task: %v\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create cleaning task", "details": err.Error()})
			return
		}
	}

	// Update booking status to 'cleaning'
	if err := s.Db.Model(&models.Book{}).Where("book_id = ?", bookingId).
		Update("status", "cleaning").Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update booking status"})
		return
	}

	if err := s.Db.Model(&models.Room{}).Where("room_id = ?", payload.RoomId).
		Update("status", "cleaning").Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update booking status"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cleaners assigned successfully",
		"count":   len(payload.CleanerIds),
	})
}

// Assign cleaner
func (s *Server) CheckinStatus(ctx *gin.Context) {
	bookingId := ctx.Param("id")

	// First get the booking to retrieve room_id
	var book models.Book
	if err := s.Db.Where("book_id = ?", bookingId).First(&book).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	// Update booking status to check-in
	if err := s.Db.Model(&book).Update("status", "check-in").Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update booking status"})
		return
	}

	// Update room status to occupied
	if err := s.Db.Model(&models.Room{}).Where("room_id = ?", book.RoomId).
		Update("status", "occupied").Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update room status"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Checking in...."})
}

// Assign cleaner
func (s *Server) CheckOut(ctx *gin.Context) {
	bookingId := ctx.Param("id")

	// First get the booking to retrieve room_id
	var book models.Book
	if err := s.Db.Where("book_id = ?", bookingId).First(&book).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	// Update booking status to check-out
	if err := s.Db.Model(&book).Update("status", "check-out").Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update booking status"})
		return
	}

	// Update room status to avail
	if err := s.Db.Model(&models.Room{}).Where("room_id = ?", book.RoomId).
		Update("status", "available").Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update room status"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Checking in...."})
}

// Assign cleaner
// Assign cleaner
func (s *Server) Completed(ctx *gin.Context) {

	var payload struct {
		RoomId string `json:"room_id"`
	}
	bookingId := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload", "details": err.Error()})
		return
	}

	if payload.RoomId == "" {
		ctx.JSON(400, gin.H{"error": "room_id is required"})
		return
	}

	if err := s.Db.Model(&models.Book{}).Where("book_id = ?", bookingId).
		Update("status", "completed").Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update booking status"})
		return
	}

	if err := s.Db.Model(&models.Room{}).Where("room_id = ?", payload.RoomId).
		Update("status", "available").Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update room status"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Booking completed and room is now available"})
}

func GenerateCleaningID(db *gorm.DB) (string, error) {
	var lastID string

	err := db.
		Model(&models.CleaningTask{}).
		Select("cleaner_id").
		Order("cleaner_id DESC").
		Limit(1).
		Scan(&lastID).Error

	if err != nil {
		return "", err
	}

	nextNumber := 1

	if lastID != "" {
		fmt.Sscanf(lastID, "CLEANTASK-%d", &nextNumber)
		nextNumber++
	}

	return fmt.Sprintf("CLEANTASK-%03d", nextNumber), nil
}
