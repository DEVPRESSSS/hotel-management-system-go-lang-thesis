package controllers

import (
	"HMS-GO/internal/models"
	"HMS-GO/internal/models/dto"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
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
