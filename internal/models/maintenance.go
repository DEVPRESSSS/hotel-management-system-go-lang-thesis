package models

import "time"

type Maintenance struct {
	Id        string    `gorm:"column:id;type:varchar(30);primaryKey" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(50);not null;uniqueIndex" json:"Name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}
