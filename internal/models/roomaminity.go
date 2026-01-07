package models

//Define database table
type RoomAminity struct {
	RoomId    string `gorm:"column:room_id;type:varchar(30);primaryKey"`
	AmenityId string `gorm:"column:amenity_id;type:varchar(30);primaryKey"`

	Room    Room    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Amenity Aminity `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
