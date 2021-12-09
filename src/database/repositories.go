package database

import (
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(description string, amount int32) uint
}

type DefaultTransactionRepository struct {
	db *gorm.DB
}

func (transactionRepo DefaultTransactionRepository) Create(description string, amount int32) uint {
	transaction := Transaction{Description: description, Amounts: amount, Created: time.Now()}
	transactionRepo.db.Create(transaction)
	return transaction.Id
}

// func (transactionRepo DefaultTransactionRepository) Find() uint {
// 	transaction := Transaction{Description: description, Amounts: amount, Created: time.Now()}
// 	transactionRepo.db.Create(transaction)
// 	return transaction.Id
// }

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return DefaultTransactionRepository{db: db}
}
