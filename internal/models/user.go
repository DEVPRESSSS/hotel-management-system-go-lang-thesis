package models

import "time"

type User struct {
	UserId    string    `gorm:"column:user_id;type:varchar(30);primaryKey" json:"userid"`
	Username  string    `gorm:"column:username;type:varchar(50);not null;uniqueIndex" json:"username"`
	FullName  string    `gorm:"column:full_name;type:varchar(100);not null" json:"fullname"`
	Email     string    `gorm:"column:email;type:varchar(100);not null;uniqueIndex" json:"email"`
	Password  string    `gorm:"column:password;type:varchar(255);not null" json:"password"`
	Locked    bool      `gorm:"column:locked;default:false" json:"locked"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	RoleId    string    `gorm:"column:role_id;type:varchar(30);not null" json:"roleid"`
	Role Role `gorm:"foreignKey:RoleId;references:RoleId" json:"role"`
}
