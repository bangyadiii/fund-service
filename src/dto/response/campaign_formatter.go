package response

import "backend-crowdfunding/src/model"

type CampaignResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	Slug             string `json:"slug"`
	ImageUrl         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	UserID           string `json:"user_id"`
}

func FormatCampaignCollections(campaigns []*model.Campaign) []CampaignResponse {
	var campaignsFormatted []CampaignResponse

	for _, data := range campaigns {
		campaign := FormatCampaign(*data)
		campaignsFormatted = append(campaignsFormatted, campaign)
	}

	return campaignsFormatted
}

func FormatCampaign(campaign model.Campaign) CampaignResponse {
	formatted := CampaignResponse{}

	formatted.ID = campaign.ID
	formatted.Name = campaign.Name
	formatted.ShortDescription = campaign.ShortDescription
	formatted.Description = campaign.Description
	formatted.Slug = campaign.Slug
	formatted.UserID = campaign.UserID
	formatted.ImageUrl = ""

	if len(campaign.CampaignImages) > 0 {
		formatted.ImageUrl = campaign.CampaignImages[0].ImageName
	}

	return formatted
}

type CampaignDetailFormatter struct {
	ID               string                `json:"id"`
	UserID           string                `json:"user_id"`
	Name             string                `json:"name"`
	ShortDescription string                `json:"short_description"`
	Description      string                `json:"description"`
	Perks            string                `json:"perks"`
	BackerCount      int                   `json:"backer_count"`
	GoalAmount       int                   `json:"goal_amount"`
	CurrentAmount    int                   `json:"current_amount"`
	Slug             string                `json:"slug"`
	CampaignImages   []model.CampaignImage `json:"campaign_images"`
}

func FormatCampaignDetail(campaign model.Campaign) CampaignDetailFormatter {
	formatted := CampaignDetailFormatter{
		ID:               campaign.ID,
		Name:             campaign.Name,
		UserID:           campaign.UserID,
		ShortDescription: campaign.ShortDescription,
		Description:      campaign.Description,
		Perks:            campaign.Perks,
		BackerCount:      campaign.BackerCount,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		Slug:             campaign.Slug,
		CampaignImages:   campaign.CampaignImages,
	}

	return formatted
}

type CampaignImageFormatter struct {
	CampaignID uint   `json:"campaign_id"`
	ImageName  string `json:"image_name"`
	IsPrimary  string `json:"is_primary"`
}
