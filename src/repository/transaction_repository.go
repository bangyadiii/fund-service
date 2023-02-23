package repository

import (
	"backend-crowdfunding/database"
	model "backend-crowdfunding/src/model"
)

type TransactionRepository interface {
	GetTransactionByCampaignID(campaignID uint) ([]model.Transaction, error)
	CreateTransaction(transaction model.Transaction) (model.Transaction, error)
}

type transactionRepository struct {
	db *database.DB
}

func NewTransactionRepository(db *database.DB) *transactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) GetTransactionByCampaignID(campaignID uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Model(model.Transaction{}).Where("campaign_id = ?", campaignID).Find(&transactions).Error

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *transactionRepository) CreateTransaction(transaction model.Transaction) (model.Transaction, error) {
	trx := r.db.Create(&transaction)
	if trx.Error != nil {
		return model.Transaction{}, trx.Error
	}

	return transaction, nil
}
