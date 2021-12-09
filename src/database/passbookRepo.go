package database

import (
	"time"

	"gorm.io/gorm"
)

type PassbookRepository interface {
	Find(owner string) (*Passbook, error)
	Create(name, owner string, totalAmount int32) uint
	CreateWithTrx(db *gorm.DB, name, owner string, totalAmount int32) (*uint, error)
	Update(passbook *Passbook) uint
	UpdateWithTrx(db *gorm.DB, passbook *Passbook) (*uint, error)
}

type DefaultPassbookRepository struct {
	db *gorm.DB
}

func NewPassbookRepository(db *gorm.DB) PassbookRepository {
	return DefaultPassbookRepository{db: db}
}

func (passbookRepository DefaultPassbookRepository) Find(owner string) (*Passbook, error) {
	var passbook Passbook
	result := passbookRepository.db.First(&passbook, "owner = ?", owner)
	if result.Error != nil {
		return nil, result.Error
	}
	return &passbook, nil
}

func (passbookRepository DefaultPassbookRepository) Create(name, owner string, totalAmount int32) uint {
	passbook := Passbook{Created: time.Now(), Updated: time.Now(), Name: name, TotalAmount: totalAmount}
	passbookRepository.db.Create(passbook)
	return passbook.Id
}

func (passbookRepository DefaultPassbookRepository) CreateWithTrx(db *gorm.DB, name, owner string, totalAmount int32) (*uint, error) {
	passbook := Passbook{Created: time.Now(), Updated: time.Now(), Name: name, TotalAmount: totalAmount}
	result := db.Create(passbook)
	if result.Error != nil {
		return nil, result.Error
	}
	return &passbook.Id, nil
}

func (passbookRepository DefaultPassbookRepository) Update(passbook *Passbook) uint {
	passbookRepository.db.Save(&passbook)
	return passbook.Id
}

func (passbookRepository DefaultPassbookRepository) UpdateWithTrx(db *gorm.DB, passbook *Passbook) (*uint, error) {
	result := db.Save(&passbook)
	if result.Error != nil {
		return nil, result.Error
	}
	return &passbook.Id, nil
}
