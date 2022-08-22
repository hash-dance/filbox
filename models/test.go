package models

import (
	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/types"
)

// Test Test struct
type Test struct {
	types.BaseModel
	A string  // 必填项
	B *string `gorm:"default:''"`  // 非必填
	C *bool   `gorm:"default:'0'"` // 非必填
	D *int    `gorm:"default:'0'"`
}

func CreateTest(test *Test) error {
	return mysql.GetClient().Create(test).Error
}

func UpdateTestByID(test *Test) error {
	return mysql.GetClient().Model(test).Updates(test).Error
}
