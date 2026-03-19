package dto

import (
	"HMS-GO/internal/models"
	"time"
)

// Book - Main booking record
type ReservationVM struct {
	BookId string `json:"book_id"`
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Email   string `json:"email"`
	Contact   string `json:"contact"`

	// Room Information
	RoomId     string `json:"room_id"`
	RoomNumber string `json:"room_number"`
	RoomType   string `json:"room_type"`

	// Booking Dates
	CheckInDate  time.Time `json:"check_in_date"`
	CheckOutDate time.Time `json:"check_out_date"`

	// Guest Information
	NumGuests int                   `json:"num_guests"`
	Guests    []models.BookingGuest `json:"guests,omitempty"`

	// Pricing
	TotalPrice    float64 `json:"total_price"`
	PricePerNight float64 `json:"price_per_night"`

	// Status
	Status string `json:"status"`

	// Payment
	PaymentStatus string `json:"payment_status"`
	// Special Requests
	SpecialRequests string `json:"special_requests"`

	// Timestamps
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CancelledAt *time.Time `json:"cancelled_at"`
}
