package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Room struct {
	RoomId     string          `json:"roomid" gorm:"primaryKey;size:36;not null"`
	RoomNumber string          `json:"roomnumber" gorm:"size:50;uniqueIndex;not null"`
	RoomTypeId string          `json:"roomtypeid" gorm:"size:30;not null"`
	RoomType   RoomType        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	FloorId    string          `gorm:"column:floor_id;type:varchar(30);not null" json:"floorid"`
	Floor      Floor           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Capacity   string          `json:"capacity"`
	Price      decimal.Decimal `json:"price" gorm:"type:decimal(10,2)"`
	Status     string          `json:"status" gorm:"size:20;default:available"`
	CreatedAt  time.Time       `json:"created_at"`

	Amenities []Amenity `json:"amenities" gorm:"many2many:room_amenities;foreignKey:RoomId;joinForeignKey:RoomId;References:AmenityId;joinReferences:AmenityId"`
}
