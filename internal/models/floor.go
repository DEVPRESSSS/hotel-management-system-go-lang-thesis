package models

type Floor struct {
	FloorId   string `gorm:"column:floor_id;type:varchar(30);primaryKey" json:"floorid"`
	FloorName string `gorm:"column:floor_name;type:varchar(50);not null;uniqueIndex" json:"floorname"`
}
