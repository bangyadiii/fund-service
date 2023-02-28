package model

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID         string `gorm:"primaryKey"`
	CampaignID string
	UserID     string
	Amount     uint
	Status     string
	Code       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
