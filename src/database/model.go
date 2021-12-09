package database

import (
	"time"
)

type Transaction struct {
	Id          uint
	Description string
	Amounts     int32
	Created     time.Time
}

type Passbook struct {
	Id          uint
	TotalAmount int32
	Name        string
	Owner       string
	Created     time.Time
	Updated     time.Time
}
