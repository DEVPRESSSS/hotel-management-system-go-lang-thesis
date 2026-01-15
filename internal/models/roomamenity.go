package models

//Define database table
type RoomAmenity struct {
	RoomId    string `gorm:"column:room_id;type:varchar(30);primaryKey"`
	AmenityId string `gorm:"column:amenity_id;type:varchar(30);primaryKey"`

	Room    Room    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Amenity Amenity `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
