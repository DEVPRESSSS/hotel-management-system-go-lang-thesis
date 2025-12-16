package models

type Role struct {
	RoleId   string `gorm:"size:30;not null;primary_key" json:"roleid"`
	RoleName string `gorm:"size:50;not null;unique" json:"rolename"`
}
