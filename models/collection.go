package models

import "github.com/google/uuid"

type Collection struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	UserID          uuid.UUID `json:"user_id" gorm:"type:uuid;not null;"`
	ContractAddress string    `json:"contract_address" gorm:"not null;"`
	CreatedAt       int64     `json:"-" gorm:"autoCreateTime:milli;"`
	UpdatedAt       int64     `json:"-" gorm:"autoUpdateTime:milli;"`
}
