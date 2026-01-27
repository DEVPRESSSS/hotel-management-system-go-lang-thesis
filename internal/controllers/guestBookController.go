package controllers

import (
	"HMS-GO/internal/models"
	"HMS-GO/internal/models/dto"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func (s *Server) RoomSelected(ctx *gin.Context) {
	roomId := ctx.Param("roomid")
	var room models.Room
	if err := s.Db.
		Preload("RoomType").
		Preload("Floor").
		Preload("Amenities").
		Where("room_id = ?", &roomId).
		First(&room).Error; err != nil {

		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching data!!!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"room": room})
}

// Calculate boooking price (pre-booking)
func (s *Server) CalculateBookingPrice(ctx *gin.Context) {

	var req dto.PriceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Get room price
	var room models.Room
	if err := s.Db.Where("room_id = ?", req.RoomID).First(&room).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	// Parse dates
	checkIn, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid check-in date"})
		return
	}

	checkOut, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid check-out date"})
		return
	}

	if !checkOut.After(checkIn) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "checkout must be after checkin"})
		return
	}

	nights := int(checkOut.Sub(checkIn).Hours() / 24)

	nightsDecimal := decimal.NewFromInt(int64(nights))
	numberGuest := decimal.NewFromInt(int64(req.Guest))
	total := nightsDecimal.Mul(room.Price).Mul(numberGuest)
	fmt.Print(total)
	ctx.JSON(http.StatusOK, gin.H{
		"price_per_night": room.Price,
		"nights":          nights,
		"total":           total,
	})
}

// Submit booking
func (s *Server) ConfirmBooking(ctx *gin.Context) {

	userId, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var book models.Book

	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate booking ID
	bookingID, err := GenerateBookingID(s.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate booking ID"})
		return
	}

	// Assign server-controlled values
	book.BookId = bookingID
	book.UserId = userId.(string)

	fmt.Print("This is the data of the guest data:", book.Guests)
	// Assign BookId to each guest
	for i := range book.Guests {
		book.Guests[i].Id = fmt.Sprintf("BKGUEST-%03d", i+1)
		book.Guests[i].BookId = bookingID
		book.Guests[i].GuestNumber = i + 1
	}

	// // Save booking + guests (GORM handles cascade)
	if err := s.Db.Create(&book).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Booking confirmation has been sent to your email!"})
}

// Generate auto IncrementId
func GenerateBookingID(db *gorm.DB) (string, error) {
	var lastID string

	err := db.
		Model(&models.Book{}).
		Select("book_id").
		Order("book_id DESC").
		Limit(1).
		Scan(&lastID).Error

	if err != nil {
		return "", err
	}

	nextNumber := 1

	if lastID != "" {
		fmt.Sscanf(lastID, "BOOKING-%d", &nextNumber)
		nextNumber++
	}

	return fmt.Sprintf("BOOKING-%03d", nextNumber), nil
}
