package campaign

import (
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	gorm.Model
	ID               int
	UserID           int 
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	CampaignImages 	 []CampaignImage
	CreatedAt        time.Time
	UpdatedAt		 time.Time
}	


type CampaignImage struct{
	gorm.Model
	ID				int
	CampaignID 		int
	ImageName 		string
	IsPrimary 		int
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
}