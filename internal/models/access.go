package models

type Access struct {
	AccessId   string `gorm:"column:access_id;type:varchar(30);primaryKey" json:"accessid"`
	AccessName string `gorm:"column:access_name;type:varchar(50);not null;uniqueIndex" json:"accessname"`
}
