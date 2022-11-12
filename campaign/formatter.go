package campaign

type CampaignFormatter struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	Slug             string `json:"slug"`
	ImageUrl         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	UserID           uint   `json:"user_id"`
}

func FormatCampaignCollections(campaigns []Campaign) []CampaignFormatter {
	var campaignsFormatted []CampaignFormatter

	for _, data := range campaigns {
		campaign := FormatCampaign(data)
		campaignsFormatted = append(campaignsFormatted, campaign)
	}

	return campaignsFormatted
}

func FormatCampaign(campaign Campaign) CampaignFormatter {
	formatted := CampaignFormatter{}

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
	ID               uint            `json:"id"`
	UserID           uint            `json:"user_id"`
	Name             string          `json:"name"`
	ShortDescription string          `json:"short_description"`
	Description      string          `json:"description"`
	Perks            string          `json:"perks"`
	BackerCount      int             `json:"backer_count"`
	GoalAmount       int             `json:"goal_amount"`
	CurrentAmount    int             `json:"current_amount"`
	Slug             string          `json:"slug"`
	CampaignImages   []CampaignImage `json:"campaign_images"`
}

func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {
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
