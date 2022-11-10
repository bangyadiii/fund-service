package campaign

import (
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	gorm.Model
	ID               uint32
	UserID           uint32
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	CampaignImages   []CampaignImage
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CampaignImage struct {
	gorm.Model
	ID         uint32
	CampaignID uint32 `gorm:"constraint:onDelete:CASCADE"`
	ImageName  string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
