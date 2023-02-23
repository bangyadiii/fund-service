package request

import (
	"backend-crowdfunding/src/model"
)

type CreateCampaignInput struct {
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description" binding:"required"`
	Perks            string `json:"perks"`
	BackerCount      uint   `json:"backer_count" binding:"required,numeric"`
	GoalAmount       uint   `json:"goal_amount" binding:"required,numeric"`
	User             model.User
}

type GetCampaignByIDInput struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateCampaignInput struct {
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description" `
	Perks            string `json:"perks"`
	BackerCount      uint   `json:"backer_count" `
	GoalAmount       uint   `json:"goal_amount"`
	CurrentAmount    uint   `json:"current_amount"`
	User             model.User
}

type UploadCampaignImageInput struct {
	CampaignID string `form:"campaign_id" binding:"required"`
	IsPrimary  bool   `form:"is_primary"`
	ImageName  string
	User       model.User
}
