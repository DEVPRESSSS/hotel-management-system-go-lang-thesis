package controllers

import (
	"HMS-GO/internal/models"
	"HMS-GO/internal/models/dto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) FetchCalendar(ctx *gin.Context) {
	var books []models.Book

	roomId := ctx.Param("room_id")

	// Only fetch future/active bookings for calendar display
	if err := s.Db.
		Where("check_out_date >= ? AND room_id = ?", time.Now(), roomId).
		Order("check_in_date ASC").
		Find(&books).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch bookings",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"books": books,
	})
}

func (s *Server) GetBookingInfo(ctx *gin.Context) {

	bookId := ctx.Param("bookId")
	var reservation models.Book
	if err := s.Db.
		Preload("User").
		Preload("Guests").
		Where("book_id", bookId).
		First(&reservation).Error; err != nil {
		ctx.JSON(http.StatusNotFound, "Failed to fetch reservation record")
		return
	}

	bookingDetails := dto.BookDetails{
		BookId:        reservation.BookId,
		RoomNumber:    reservation.RoomNumber,
		Fullname:      reservation.User.FullName,
		RoomType:      reservation.RoomType,
		CheckInDate:   reservation.CheckInDate,
		CheckOutDate:  reservation.CheckOutDate,
		NumGuests:     reservation.NumGuests,
		Guests:        reservation.Guests,
		TotalPrice:    reservation.TotalPrice,
		Status:        reservation.Status,
		PaymentStatus: reservation.PaymentStatus,
	}
	ctx.JSON(http.StatusOK, gin.H{"booking_details": bookingDetails})
}
