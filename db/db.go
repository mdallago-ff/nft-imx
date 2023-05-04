package db

import (
	"nft/models"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func (d *DB) UpdateUser(user *models.User) error {
	//save database
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(&user).Error; err != nil {
			return err
		}

		return nil
	})
}

func (d *DB) GetUser(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := d.db.First(&user, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (d *DB) GetUserByMail(mail string) (*models.User, error) {
	var user models.User
	if err := d.db.Where("mail = ?", mail).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (d *DB) CreateCollection(collection *models.Collection) error {
	//save database
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&collection).Error; err != nil {
			return err
		}

		return nil
	})
}

func (d *DB) GetCollection(id uuid.UUID) (*models.Collection, error) {
	var collection models.Collection
	if err := d.db.First(&collection, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, err
		}
	}

	return &collection, nil
}

func (d *DB) CreateToken(token *models.Token) error {
	//save database
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&token).Error; err != nil {
			return err
		}

		return nil
	})
}

func (d *DB) GetToken(id uuid.UUID) (*models.Token, error) {
	var token models.Token
	if err := d.db.First(&token, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, err
		}
	}

	return &token, nil
}
