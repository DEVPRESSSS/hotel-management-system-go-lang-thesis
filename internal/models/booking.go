package models

import "time"

// Book - Main booking record
type Book struct {
	BookId string `gorm:"column:book_id;type:varchar(36);primaryKey" json:"book_id"`
	UserId string `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	User   User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	// Room Information
	RoomId     string `gorm:"column:room_id;type:varchar(36);not null" json:"room_id"`
	RoomNumber string `gorm:"column:room_number;type:varchar(50)" json:"room_number"`
	RoomType   string `gorm:"column:room_type;type:varchar(50)" json:"room_type"`

	// Booking Dates
	CheckInDate  time.Time `gorm:"column:check_in_date;not null" json:"check_in_date"`
	CheckOutDate time.Time `gorm:"column:check_out_date;not null" json:"check_out_date"`

	// Guest Information
	NumGuests int            `gorm:"column:num_guests;default:1" json:"num_guests"`
	Guests    []BookingGuest `gorm:"foreignKey:BookId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"guests,omitempty"`

	// Pricing
	TotalPrice    float64 `gorm:"column:total_price;type:decimal(10,2)" json:"total_price"`
	PricePerNight float64 `gorm:"column:price_per_night;type:decimal(10,2)" json:"price_per_night"`

	// Status
	Status string `gorm:"column:status;type:varchar(20);default:'pending'" json:"status"`
	// pending, confirmed, checked_in, checked_out, cancelled

	// Payment
	PaymentStatus string `gorm:"column:payment_status;type:varchar(20);default:'unpaid'" json:"payment_status"`
	// unpaid, partial, paid, refunded

	// Special Requests
	SpecialRequests string `gorm:"column:special_requests;type:text" json:"special_requests"`

	// Timestamps
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	CancelledAt *time.Time `gorm:"column:cancelled_at" json:"cancelled_at,omitempty"`
}
