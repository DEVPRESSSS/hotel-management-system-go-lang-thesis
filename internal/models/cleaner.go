package models

import "time"

type Cleaner struct {
	Id        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(50);not null" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}
