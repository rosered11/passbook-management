package services

import (
	"rosered/passbook-management/src/database"
	"rosered/passbook-management/src/dto"
)

type PassbookService interface {
	Add(request dto.PassbookRequest) error
}

type DefaultPassbookService struct {
	unitofwork database.UnitOfWork
}

func (passbookService DefaultPassbookService) Add(request dto.PassbookRequest) error {
	trx := passbookService.unitofwork.DB.Begin()
	_, err := passbookService.unitofwork.TransactionRepo.CreateWithTrx(trx, request.Description, request.Amount)
	if err != nil {
		trx.Rollback()
		return err
	}
	passbook, _ := passbookService.unitofwork.PassbookRepo.Find(request.Owner)
	if passbook == nil {
		_, err = passbookService.unitofwork.PassbookRepo.CreateWithTrx(trx, "Default", request.Owner, request.Amount)
		if err != nil {
			trx.Rollback()
			return err
		}
	} else {
		passbook.TotalAmount += request.Amount
		_, err = passbookService.unitofwork.PassbookRepo.UpdateWithTrx(trx, passbook)
		if err != nil {
			trx.Rollback()
			return err
		}
	}
	trx.Commit()
	return nil
}

func NewPassbookService(unitofwork database.UnitOfWork) PassbookService {
	return DefaultPassbookService{unitofwork: unitofwork}
}
