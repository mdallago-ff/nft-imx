package models

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Mail      string    `json:"email" gorm:"not null;unique"`
	ApiKey    string    `json:"api_key" gorm:"not null;"`
	Private   string    `json:"-" gorm:"not null;"`
	Public    string    `json:"public" gorm:"not null;"`
	Address   string    `json:"address" gorm:"not null;"`
	StarkKey  string    `json:"-" gorm:"null;"`
	CreatedAt int64     `json:"-" gorm:"autoCreateTime:milli;"`
	UpdatedAt int64     `json:"-" gorm:"autoUpdateTime:milli;"`
}
