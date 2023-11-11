package request

type TransactionByCampaignInput struct {
	CampaignID string `json:"campaign_id" validate:"required"`
}

type CreateTransactionInput struct {
	CampaignID string `json:"campaign_id" validate:"required"`
	UserID     string
	Amount     uint   `json:"amount" validate:"required"`
	Status     string `json:"status" validate:"required"`
	Code       string `json:"code" validate:"required"`
}
