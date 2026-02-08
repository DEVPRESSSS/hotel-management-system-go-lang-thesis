package models

import "time"

type HistoryLog struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	EntityType  string `gorm:"size:100;not null"`
	EntityID    string `gorm:"size:100;not null"`
	Action      string `gorm:"size:50;not null"`
	Description string `gorm:"type:text"`
	OldValue    string `gorm:"type:json"`
	NewValue    string `gorm:"type:json"`
	PerformedBy *uint64
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
