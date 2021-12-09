package database

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Description string
	Amounts     int32
}

type Passbook struct {
	gorm.Model
	TotalAmount int32
	Name        string
	Owner       string
}
