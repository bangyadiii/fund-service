package request

type TransactionByCampaignInput struct {
	CampaignID uint `json:"campaign_id" binding:"required"`
}

type CreateTransactionInput struct {
	CampaignID uint `json:"campaign_id" binding:"required"`
	UserID     uint
	Amount     uint   `json:"amount" binding:"required"`
	Status     string `json:"status" binding:"required"`
	Code       string `json:"code" binding:"required"`
}
