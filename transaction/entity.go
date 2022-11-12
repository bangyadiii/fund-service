package transaction

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID         uint
	CampaignID uint
	UserID     uint
	Amount     uint
	Status     string
	Code       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
