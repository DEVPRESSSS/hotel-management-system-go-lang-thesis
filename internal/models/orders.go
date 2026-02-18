package models

import "time"

type Orders struct {
	OrderId string `gorm:"column:order_id;type:varchar(36);primaryKey" json:"order_id"`
	UserId  string `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	BookId string `gorm:"column:book_id;type:varchar(36);not null" json:"book_id"`
	Book   Book   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	// Room Information
	ProductName string `gorm:"column:product_name;type:varchar(50)" json:"product_name"`
	SessionId   string `gorm:"column:session_id;type:varchar(50)" json:"session_id"`
	Qty         int    `gorm:"column:qty;type:int" json:"qty"`
	// Pricing
	TotalPrice float64 `gorm:"column:total_price;type:decimal(10,2)" json:"total_price"`
	Price      float64 `gorm:"column:price_per_night;type:decimal(10,2)" json:"price_per_night"`

	// Payment
	PaymentStatus string `gorm:"column:payment_status;type:varchar(20);default:'unpaid'" json:"payment_status"`

	// Timestamps
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}
