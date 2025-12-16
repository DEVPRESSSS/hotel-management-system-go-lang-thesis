package models

type User struct {
	UserId    string `gorm:"size:30;not null;primaryKey" json:"userid"`
	Username  string `gorm:"size:50;not null;unique" json:"username"`
	FullName  string `gorm:"size:50;not null;unique" json:"fullname"`
	Email     string `gorm:"size:50;not null;unique" json:"email"`
	Password  string `gorm:"size:30;not null" json:"password"`
	Locked    bool   `gorm:"default:false" json:"locked"`
	CreatedAt string `gorm:"autoCreateTime" json:"created_at"`

	RoleId string `gorm:"size:30;not null" json:"roleid"`
	Role   Role   `gorm:"foreignKey:RoleId;references:RoleId" json:"role"`
}
