package request

type TransactionByCampaignInput struct {
	CampaignID string `json:"campaign_id" binding:"required"`
}

type CreateTransactionInput struct {
	CampaignID string `json:"campaign_id" binding:"required"`
	UserID     string
	Amount     uint   `json:"amount" binding:"required"`
	Status     string `json:"status" binding:"required"`
	Code       string `json:"code" binding:"required"`
}
