package database

import (
	"gorm.io/gorm"
)

type PassbookRepository interface {
	Find(owner string) (*Passbook, error)
	FindWithDate(owner string, currentDate string) (*Passbook, error)
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

func (passbookRepository DefaultPassbookRepository) FindWithDate(owner string, currentDate string) (*Passbook, error) {
	var passbook Passbook
	result := passbookRepository.db.First(&passbook, "owner = ? AND updated_at >= ?", owner, currentDate)
	if result.Error != nil {
		return nil, result.Error
	}
	return &passbook, nil
}

func (passbookRepository DefaultPassbookRepository) Create(name, owner string, totalAmount int32) uint {
	passbook := Passbook{Name: name, TotalAmount: totalAmount, Owner: owner}
	passbookRepository.db.Create(&passbook)
	return passbook.ID
}

func (passbookRepository DefaultPassbookRepository) CreateWithTrx(db *gorm.DB, name, owner string, totalAmount int32) (*uint, error) {
	passbook := Passbook{Name: name, TotalAmount: totalAmount, Owner: owner}
	result := db.Create(&passbook)
	if result.Error != nil {
		return nil, result.Error
	}
	return &passbook.ID, nil
}

func (passbookRepository DefaultPassbookRepository) Update(passbook *Passbook) uint {
	passbookRepository.db.Save(&passbook)
	return passbook.ID
}

func (passbookRepository DefaultPassbookRepository) UpdateWithTrx(db *gorm.DB, passbook *Passbook) (*uint, error) {
	result := db.Save(&passbook)
	if result.Error != nil {
		return nil, result.Error
	}
	return &passbook.ID, nil
}
