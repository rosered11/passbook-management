package database

import (
	"gorm.io/gorm"
)

const DEFAULT_PASSBOOK_NAME string = "passbook"

type PassbookRepository interface {
	FindWithOwner(db *gorm.DB, passbook interface{}, owner string) error
	FindWithOwnerAndDate(db *gorm.DB, passbook interface{}, owner string, currentDate string) error
	FindWithDate(db *gorm.DB, datenow string) (*[]Passbook, error)
	Create(db *gorm.DB, passbook *Passbook) (*uint, error)
	Update(db *gorm.DB, passbook *Passbook) (*uint, error)
}

type DefaultPassbookRepository struct {
	repository Repository
}

func NewPassbookRepository(repository Repository) PassbookRepository {
	return DefaultPassbookRepository{repository: repository}
}

func (passbookRepository DefaultPassbookRepository) FindWithOwner(db *gorm.DB, passbook interface{}, owner string) error {
	passbookRepository.repository.Find(db, passbook, "owner = ? AND name = ?", owner, DEFAULT_PASSBOOK_NAME)
	return nil
}

func (passbookRepository DefaultPassbookRepository) FindWithDate(db *gorm.DB, datenow string) (*[]Passbook, error) {
	var passbooks []Passbook
	_, err := passbookRepository.repository.Find(db, &passbooks, "updated_at >= ?", datenow)
	if err != nil {
		return nil, err
	}
	return &passbooks, nil
}

func (passbookRepository DefaultPassbookRepository) FindWithOwnerAndDate(db *gorm.DB, passbook interface{}, owner string, currentDate string) error {
	passbookRepository.repository.Find(db, passbook, "owner = ? AND name = ? AND updated_at >= ?", owner, DEFAULT_PASSBOOK_NAME, currentDate)
	return nil
}

func (passbookRepository DefaultPassbookRepository) Create(db *gorm.DB, passbook *Passbook) (*uint, error) {
	passbook.Name = DEFAULT_PASSBOOK_NAME
	result := db.Create(passbook)
	if result.Error != nil {
		return nil, result.Error
	}
	return &passbook.ID, nil
}

func (passbookRepository DefaultPassbookRepository) Update(db *gorm.DB, passbook *Passbook) (*uint, error) {
	result := db.Save(passbook)
	if result.Error != nil {
		return nil, result.Error
	}
	return &passbook.ID, nil
}
