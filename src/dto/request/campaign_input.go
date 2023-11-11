package request

import (
	"backend-crowdfunding/src/model"
)

type CreateCampaignInput struct {
	Name             string `json:"name" validate:"required"`
	ShortDescription string `json:"short_description" validate:"required, max=255"`
	Description      string `json:"description" validate:"required, max=900"`
	Perks            string `json:"perks"`
	BackerCount      uint   `json:"backer_count" validate:"required,numeric"`
	GoalAmount       uint   `json:"goal_amount" validate:"required,numeric"`
}

type GetCampaignByIDInput struct {
	ID string `uri:"id" validate:"required"`
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
	CampaignID string `form:"campaign_id" validate:"required"`
	IsPrimary  bool   `form:"is_primary"`
	ImageName  string
	User       model.User
}

type CampaignsWithPaginationParam struct {
	Name   string `json:"name,omitempty;query:name"`
	UserID string `json:"user_id,omitempty"`
	PaginationParam
}
