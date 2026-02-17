package models

type FoodCategory struct {
	FoodCategoryId string `gorm:"column:food_category_id;type:varchar(30);primaryKey" json:"foodcategoryId"`
	Name           string `gorm:"column:name;type:varchar(50);not null;uniqueIndex" json:"name"`
	Time           string `gorm:"column:time;type:varchar(50);not null;" json:"time"`
}
