package models

type Aminity struct {
	AminityId   string `gorm:"column:aminity_id;type:varchar(30);primaryKey" json:"aminityid"`
	AminityName string `gorm:"column:aminity_name;type:varchar(50);not null;uniqueIndex" json:"aminityname"`
}
