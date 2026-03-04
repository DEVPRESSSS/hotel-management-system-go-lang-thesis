package dto

import (
	"HMS-GO/internal/models"
	"time"
)

type BookDetails struct {
	BookId   string `json:"book_id"`
	Fullname string `json:"fullname"`
	// Room Information
	RoomNumber string `json:"room_number"`
	RoomType   string `json:"room_type"`
	// Booking Dates
	CheckInDate  time.Time `json:"check_in_date"`
	CheckOutDate time.Time `json:"check_out_date"`
	// Guest Information
	NumGuests int                   `json:"num_guests"`
	Guests    []models.BookingGuest `json:"guests"`
	// Pricing
	TotalPrice float64 `json:"total_price"`
	// Status
	Status string `json:"status"`
	// Payment
	PaymentStatus string `json:"payment_status"`
}
