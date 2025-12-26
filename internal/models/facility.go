package models

import (
	"time"
)

type Facility struct {
	FacilityId      string    `gorm:"column:facility_id;type:varchar(30);primaryKey" json:"facility_id"`
	FacilityName    string    `gorm:"column:facility_name;type:varchar(50);not null;uniqueIndex" json:"facility_name"`
	Status          bool      `gorm:"column:status;default:false" json:"status"`
	MaintenanceDate time.Time `gorm:"type:date;column:maintenance_date" json:"maintenance_date"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}
