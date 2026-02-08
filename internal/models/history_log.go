package models

import "time"

type HistoryLog struct {
	Id          string `gorm:"column:history_id;type:varchar(30);primaryKey" json:"id"`
	EntityType  string `gorm:"size:100;not null"`
	EntityID    string `gorm:"size:100;not null"`
	Action      string `gorm:"size:50;not null"`
	Description string `gorm:"type:text"`
	OldValue    string `gorm:"type:text"`
	NewValue    string `gorm:"type:text"`
	PerformedBy string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
