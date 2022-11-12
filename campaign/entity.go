package campaign

import (
	"backend-crowdfunding/user"
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	ID               uint `gorm:"primaryKey"`
	User             user.User 
	UserID           uint
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
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

type CampaignImage struct {
	ID         uint `gorm:"primaryKey"`
	CampaignID uint
	Campaign   Campaign
	ImageName  string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
