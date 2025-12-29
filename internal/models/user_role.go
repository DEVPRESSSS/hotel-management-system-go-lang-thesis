package models

type UserRole struct {
	UserID string `gorm:"column:user_id;primaryKey"`
	RoleID string `gorm:"column:role_id;primaryKey"`

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Role Role `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
