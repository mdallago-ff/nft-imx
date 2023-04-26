package models

import "github.com/google/uuid"

type Token struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	CollectionID uuid.UUID `json:"collection_id" gorm:"type:uuid;not null;"`
	TokenID      string    `json:"token_id" gorm:"not null;"`
	CreatedAt    int64     `json:"-" gorm:"autoCreateTime:milli;"`
	UpdatedAt    int64     `json:"-" gorm:"autoUpdateTime:milli;"`
}
