package models

type Amenity struct {
	AmenityId   string `gorm:"column:amenity_id;type:varchar(30);primaryKey" json:"amenityid"`
	AmenityName string `gorm:"column:amenity_name;type:varchar(50);not null;uniqueIndex" json:"amenityname"`
}
