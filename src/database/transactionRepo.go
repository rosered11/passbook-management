package database

import (
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(db *gorm.DB, transaction *Transaction) (*uint, error)
}

type DefaultTransactionRepository struct {
	repository Repository
}

func (transactionRepo DefaultTransactionRepository) Create(db *gorm.DB, transaction *Transaction) (*uint, error) {
	result := db.Create(transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transaction.ID, nil
}

func NewTransactionRepository(repository Repository) TransactionRepository {
	return DefaultTransactionRepository{repository: repository}
}
