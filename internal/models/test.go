package models

type Test struct {
	TestId   string `gorm:"column:test_id;type:varchar(30);primaryKey" json:"testid"`
	TestName string `gorm:"column:test_name;type:varchar(50);not null;uniqueIndex" json:"testname"`
}
