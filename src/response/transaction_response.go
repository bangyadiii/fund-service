package response

import "backend-crowdfunding/src/model"

type TransactionResponse struct {
	ID         string `json:"id"`
	CampaignID string `json:"campaign_id"`
	UserID     string `json:"user_id"`
	Amount     uint   `json:"amount"`
	Status     string `json:"status"`
	Code       string `json:"code"`
}

func FormatTransaction(transaction model.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:         transaction.ID,
		CampaignID: transaction.CampaignID,
		UserID:     transaction.UserID,
		Amount:     transaction.Amount,
		Status:     transaction.Status,
		Code:       transaction.Code,
	}
}

func FormatTransactionCollections(transactions []model.Transaction) []TransactionResponse {
	var transactionsRes []TransactionResponse
	for _, transaction := range transactions {
		trx := FormatTransaction(transaction)
		transactionsRes = append(transactionsRes, trx)
	}
	return transactionsRes
}
