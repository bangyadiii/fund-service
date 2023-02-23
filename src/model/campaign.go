package model

import (
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	ID               string `gorm:"primaryKey"`
	User             User
	UserID           string
	Transactions     []Transaction
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
	ID         string `gorm:"primaryKey"`
	CampaignID string
	ImageName  string
	IsPrimary  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
