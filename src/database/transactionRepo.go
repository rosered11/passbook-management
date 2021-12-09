package database

import (
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(description string, amount int32) uint
	CreateWithTrx(db *gorm.DB, description string, amount int32) (*uint, error)
}

type DefaultTransactionRepository struct {
	db *gorm.DB
}

func (transactionRepo DefaultTransactionRepository) Create(description string, amount int32) uint {
	transaction := Transaction{Description: description, Amounts: amount}
	transactionRepo.db.Create(&transaction)
	return transaction.ID
}

func (transactionRepo DefaultTransactionRepository) CreateWithTrx(db *gorm.DB, description string, amount int32) (*uint, error) {
	transaction := Transaction{Description: description, Amounts: amount}
	result := db.Create(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transaction.ID, nil
}

// func (transactionRepo DefaultTransactionRepository) Find() uint {
// 	transaction := Transaction{Description: description, Amounts: amount, Created: time.Now()}
// 	transactionRepo.db.Create(transaction)
// 	return transaction.Id
// }

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return DefaultTransactionRepository{db: db}
}
