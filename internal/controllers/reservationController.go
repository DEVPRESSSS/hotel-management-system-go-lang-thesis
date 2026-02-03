package controllers

import (
	"HMS-GO/internal/models"
	"HMS-GO/internal/models/dto"
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

// func (s *Server) GetAllEventsReservations(ctx *gin.Context) {

// 	var books []models.Book
// 	if err := s.Db.Preload("User").Find(&books).Error; err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch bookings"})
// 		return
// 	}

// 	events := make([]dto.Calendar, 0)

// 	for _, b := range books {
// 		events = append(events, dto.Calendar{
// 			Start:      b.CheckInDate.Format("2006-01-02"),
// 			End:        b.CheckOutDate.Format("2006-01-02"),
// 			Display:    "background",
// 			Background: "#ef4444",
// 		})
// 		events = append(events, dto.Calendar{
// 			Title:     b.BookId + " " + b.RoomNumber + " (" + b.User.FullName + ")",
// 			Start:     b.CheckInDate.Format("2006-01-02"),
// 			End:       b.CheckOutDate.Format("2006-01-02"),
// 			Color:     "#ef4444",
// 			TextColor: "#ffffff",
// 		})
// 	}

//		ctx.JSON(http.StatusOK, events)
//	}
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
