package models

import "time"

type BookingGuest struct {
	Id     string `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	BookId string `gorm:"column:book_id;type:varchar(36);not null;index" json:"book_id"`

	GuestNumber int `gorm:"column:guest_number;not null" json:"guest_number"`

	// Personal Information
	FirstName   string `gorm:"column:first_name;type:varchar(100);not null" json:"firstname"`
	LastName    string `gorm:"column:last_name;type:varchar(100);not null" json:"lastname"`
	PhoneNumber string `gorm:"column:phone_number;type:varchar(20);not null" json:"phonenumber"`

	// Timestamps
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
