package models

import "time"

type CleaningTask struct {
	Id string `gorm:"type:varchar(36);primaryKey"`

	BookId string `gorm:"type:varchar(36);not null" json:"book_id"`
	Book   Book   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	RoomId string `gorm:"type:varchar(36);not null" json:"room_id"`

	CleanerId *string  `gorm:"type:varchar(36)" json:"cleaner_id"`
	Cleaner   *Cleaner `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Status string `gorm:"type:varchar(20);default:'pending'" json:"status"`

	StartedAt   *time.Time
	CompletedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
