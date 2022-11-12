package transaction

type GetTransactionByCampaignInput struct {
	CampaignID uint `json:"campaign_id" binding:"required"`
}
