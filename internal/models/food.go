package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Food struct {
	FoodId         string          `gorm:"column:food_id;type:varchar(30);primaryKey" json:"foodId"`
	Name           string          `gorm:"column:name;type:varchar(30)" json:"name"`
	Description    string          `gorm:"column:description;type:varchar(200)" json:"description"`
	Image          string          `gorm:"column:image;type:varchar(200)" json:"image"`
	FoodCategoryId string          `gorm:"column:food_category_id;type:varchar(30);not null" json:"foodcategoryid"`
	FoodCategory   FoodCategory    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Price          decimal.Decimal `json:"price" gorm:"type:decimal(10,2)"`
	Status         string          `json:"status" gorm:"size:20;default:available"`
	CreatedAt      time.Time       `json:"created_at"`
}
