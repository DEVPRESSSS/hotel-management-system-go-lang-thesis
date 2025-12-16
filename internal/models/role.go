package models

type Role struct {
	RoleId   string `gorm:"column:role_id;type:varchar(30);primaryKey" json:"roleid"`
	RoleName string `gorm:"column:role_name;type:varchar(50);not null;uniqueIndex" json:"rolename"`
}
