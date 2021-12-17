package services

import (
	"fmt"
	"rosered/passbook-management/src/database"
	"rosered/passbook-management/src/dto"
	"rosered/passbook-management/src/utilities"

	"github.com/shamaton/zeroformatter/datetimeoffset"
)

type PassbookService struct {
	unitofwork database.UnitOfWork
}

// Get Record transaction
func (passbookService PassbookService) FindTransactionOfPassbooksWithCurrentDate() (*[]database.Transaction, error) {
	now := datetimeoffset.Now()
	datenow := now.Format("2006-01-02") + " 00:00:00"
	// fmt.Printf(now.Format(time.RFC3339))
	// fmt.Printf(now.Format("2006-01-02"))
	var transaction []database.Transaction
	_, err := passbookService.unitofwork.Repository.Find(passbookService.unitofwork.DB, &transaction, "updated_at >= ?", datenow)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (passbookService PassbookService) FindTransactionOfPassbooksWithCurrentDateAndOwner(owner string) (*[]database.Transaction, error) {
	now := datetimeoffset.Now()
	datenow := now.Format("2006-01-02") + " 00:00:00"
	var transaction []database.Transaction
	_, err := passbookService.unitofwork.Repository.Find(passbookService.unitofwork.DB, &transaction, "updated_at >= ? AND owner = ?", datenow, owner)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (passbookService PassbookService) FindTransactionOfPassbooksWithDate(date string) (*[]database.Transaction, error) {
	startDate := date + " 00:00:00"
	endDate, _ := utilities.GetNextMonth(startDate)
	fmt.Printf("endate: " + endDate)
	// fmt.Printf(now.Format(time.RFC3339))
	// fmt.Printf(now.Format("2006-01-02"))
	var transaction []database.Transaction
	_, err := passbookService.unitofwork.Repository.Find(passbookService.unitofwork.DB, &transaction, "updated_at >= ? AND updated_at < ?", startDate, endDate)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (passbookService PassbookService) FindTransactionOfPassbooksWithDateAndOwner(date string, owner string) (*[]database.Transaction, error) {
	startDate := date + " 00:00:00"
	endDate, _ := utilities.GetNextMonth(startDate)
	// fmt.Printf(now.Format(time.RFC3339))
	// fmt.Printf(now.Format("2006-01-02"))
	var transaction []database.Transaction
	_, err := passbookService.unitofwork.Repository.Find(passbookService.unitofwork.DB, &transaction, "updated_at >= ? AND updated_at < ? AND owner = ?", startDate, endDate, owner)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (passbookService PassbookService) FindPassbooksWithOwner(owner string) (*[]database.Passbook, error) {
	var passbooks []database.Passbook
	err := passbookService.unitofwork.PassbookRepo.FindWithOwner(passbookService.unitofwork.DB, &passbooks, owner)
	if err != nil {
		return nil, err
	}

	return &passbooks, nil
}

func (passbookService PassbookService) FindPassbooks() (*[]database.Passbook, error) {
	var passbooks []database.Passbook
	err := passbookService.unitofwork.Repository.FindAll(passbookService.unitofwork.DB, &passbooks)
	if err != nil {
		return nil, err
	}

	return &passbooks, nil
}

func (passbookService PassbookService) Add(request *dto.PassbookRequest) error {

	trx := passbookService.unitofwork.DB.Begin()
	transaction := database.Transaction{Description: request.Description, Amounts: request.Amount, Owner: request.Owner}
	_, err := passbookService.unitofwork.TransactionRepo.Create(trx, &transaction)
	if err != nil {
		trx.Rollback()
		return err
	}

	var passbook database.Passbook
	hasRecord, err := passbookService.unitofwork.Repository.First(passbookService.unitofwork.DB, &passbook, "owner = ?", request.Owner)
	if err != nil {
		trx.Rollback()
		return err
	}
	if !hasRecord {
		passbookModel := database.Passbook{Owner: request.Owner, TotalAmount: request.Amount}
		_, err = passbookService.unitofwork.PassbookRepo.Create(trx, &passbookModel)
		if err != nil {
			trx.Rollback()
			return err
		}
	} else {
		passbook.TotalAmount += request.Amount
		_, err = passbookService.unitofwork.PassbookRepo.Update(trx, &passbook)
		if err != nil {
			trx.Rollback()
			return err
		}
	}
	trx.Commit()
	return nil
}

func NewPassbookService(unitofwork database.UnitOfWork) PassbookService {
	return PassbookService{unitofwork: unitofwork}
}
