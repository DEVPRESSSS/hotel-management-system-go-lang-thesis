package models

type RoomType struct {
	RoomTypeId   string `gorm:"column:room_typeid;type:varchar(30);primaryKey" json:"roomtypeid"`
	RoomTypeName string `gorm:"column:room_type_name;type:varchar(50);not null;uniqueIndex" json:"roomtypename"`
	Description  string `gorm:"column:description:varchar(200);not null" json:"description"`
}
