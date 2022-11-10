package campaign

type CampaignFormatter struct {
	ID               uint32          `json:"id"`
	UserID           uint32          `json:"user_id"`
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

func FormatCampaign(campaign Campaign) CampaignFormatter {
	formatted := CampaignFormatter{
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
	}

	return formatted
}

func FormatCampaignCollection(campaign []Campaign) []CampaignFormatter {

	return nil
}
