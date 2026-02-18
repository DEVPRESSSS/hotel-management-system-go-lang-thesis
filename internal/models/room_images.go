package models

type RoomImages struct {
	ImageId string `gorm:"column:image_id;type:varchar(30);primaryKey" json:"image_id"`
	RoomId  string `gorm:"column:room_id;type:varchar(30);not null" json:"room_id"`
	Room    Room   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Image   string `gorm:"column:image;type:varchar(200)" json:"image"`
}
