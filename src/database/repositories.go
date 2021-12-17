package database

import (
	"gorm.io/gorm"
)

type Repository interface {
	Find(db *gorm.DB, dest interface{}, query string, args ...interface{}) (interface{}, error)
	FindAll(db *gorm.DB, dest interface{}) error
	First(db *gorm.DB, dest interface{}, query string, args ...interface{}) (bool, error)
	Create(db *gorm.DB, dest interface{}) *gorm.DB
	Update(db *gorm.DB, dest interface{}) *gorm.DB
}

type DefaultRepository struct{}

func NewRepository() DefaultRepository {
	return DefaultRepository{}
}

func (repository DefaultRepository) Find(db *gorm.DB, dest interface{}, query string, args ...interface{}) (interface{}, error) {
	result := db.Where(query, args...).Order("id desc").Find(dest)
	if result.Error != nil {
		return nil, result.Error
	}
	return dest, nil
}

func (repository DefaultRepository) FindAll(db *gorm.DB, dest interface{}) error {
	result := db.Order("id desc").Find(dest)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repository DefaultRepository) First(db *gorm.DB, dest interface{}, query string, args ...interface{}) (bool, error) {
	result := db.Where(query, args...).Order("id desc").First(dest)
	hasRecord := false
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return hasRecord, nil
		}
		return hasRecord, result.Error
	}
	hasRecord = true
	return hasRecord, nil
}

func (repository DefaultRepository) Create(db *gorm.DB, dest interface{}) *gorm.DB {
	return db.Create(&dest)
}

func (repository DefaultRepository) Update(db *gorm.DB, dest interface{}) *gorm.DB {
	return db.Save(&dest)
}
