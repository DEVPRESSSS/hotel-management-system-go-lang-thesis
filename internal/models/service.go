package models

import "time"

type Service struct {
	ServiceId   string    `gorm:"column:service_id;type:varchar(30);primaryKey" json:"serviceid"`
	ServiceName string    `gorm:"column:service_name;type:varchar(50);not null;uniqueIndex" json:"servicename"`
	StartTime   string    `gorm:"column:start_time;type:time;not null" json:"start_time"`
	EndTime     string    `gorm:"column:end_time;type:time;not null" json:"end_time"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}
