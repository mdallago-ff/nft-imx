package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"nft/models"
)

type DB struct {
	db *gorm.DB
}

func NewDB(dsn string) (*DB, error) {
	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (d *DB) CreateUser(user *models.User) error {
	//save database
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		return nil
	})
}
