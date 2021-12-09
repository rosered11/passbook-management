package services

import (
	"rosered/passbook-management/src/database"
	"rosered/passbook-management/src/dto"

	"github.com/shamaton/zeroformatter/datetimeoffset"
)

type PassbookService interface {
	Add(request *dto.PassbookRequest) error
	FindPassbookCurrentDate(owner string) (*database.Passbook, error)
}

type DefaultPassbookService struct {
	unitofwork database.UnitOfWork
}

func (passbookService DefaultPassbookService) FindPassbookCurrentDate(owner string) (*database.Passbook, error) {
	now := datetimeoffset.Now()
	datenow := now.Format("2006-01-02") + " 00:00:00"
	// fmt.Printf(now.Format(time.RFC3339))
	// fmt.Printf(now.Format("2006-01-02"))
	result, err := passbookService.unitofwork.PassbookRepo.FindWithDate(owner, datenow)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (passbookService DefaultPassbookService) Add(request *dto.PassbookRequest) error {

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
