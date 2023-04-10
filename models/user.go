package models

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Mail      string    `json:"email" gorm:"not null;"`
	ApiKey    string    `json:"api_key" gorm:"not null;"`
	CreatedAt int64     `gorm:"autoCreateTime:milli;"`
	UpdatedAt int64     `gorm:"autoUpdateTime:milli;"`
}
