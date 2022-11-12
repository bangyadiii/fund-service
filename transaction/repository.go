package transaction

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetTransactionByCampaignID(campaignID uint) ([]Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetTransactionByCampaignID(campaignID uint) ([]Transaction, error) {
	var transactions []Transaction
	err := r.db.Model(Transaction{}).Where("campaign_id = ?", campaignID).Find(&transactions).Error

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
