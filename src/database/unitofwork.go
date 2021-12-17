package database

import (
	"gorm.io/gorm"
)

type UnitOfWork struct {
	DB              *gorm.DB
	TransactionRepo TransactionRepository
	PassbookRepo    PassbookRepository
	Repository      Repository
}

func NewUnitOfWork(db *gorm.DB, repository Repository) UnitOfWork {
	return UnitOfWork{DB: db, TransactionRepo: NewTransactionRepository(repository), PassbookRepo: NewPassbookRepository(repository), Repository: repository}
}
