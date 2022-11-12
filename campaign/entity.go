package campaign

import (
	"backend-crowdfunding/transaction"
	"backend-crowdfunding/user"
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	ID               uint `gorm:"primaryKey"`
	User             user.User
	UserID           uint
	Transactions     []transaction.Transaction
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	Slug             string `gorm:"size:256;uniqueIndex"`
	CampaignImages   []CampaignImage
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

type CampaignImage struct {
	ID         uint `gorm:"primaryKey"`
	CampaignID uint
	ImageName  string
	IsPrimary  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
