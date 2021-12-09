package database

import (
	"gorm.io/gorm"
)

type UnitOfWork struct {
	DB              *gorm.DB
	TransactionRepo TransactionRepository
	PassbookRepo    PassbookRepository
}

func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return UnitOfWork{DB: db, TransactionRepo: NewTransactionRepository(db), PassbookRepo: NewPassbookRepository(db)}
}
