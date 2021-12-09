package services

import "rosered/passbook-management/src/database"

type PassbookService interface{}

type DefaultPassbookService struct {
	transactionRepo database.TransactionRepository
}

func (passbookService DefaultPassbookService) Add() {
	passbookService.transactionRepo.Create("test", 5)
}

func NewPassbookService(transactionRepo database.TransactionRepository) PassbookService {
	return DefaultPassbookService{transactionRepo: transactionRepo}
}
