package models

type RoleAccess struct {
	RoleID   string `gorm:"column:role_id;type:varchar(30);primaryKey"`
	AccessID string `gorm:"column:access_id;type:varchar(30);primaryKey"`

	Role   Role   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Access Access `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
