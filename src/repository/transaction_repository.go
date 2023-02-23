package repository

import (
	"backend-crowdfunding/database"
	model "backend-crowdfunding/src/model"
	"backend-crowdfunding/src/util/id"
)

type TransactionRepository interface {
	GetTransactionByCampaignID(campaignID uint) ([]model.Transaction, error)
	CreateTransaction(transaction model.Transaction) (model.Transaction, error)
}

type trxRepoImpl struct {
	db          *database.DB
	idGenerator id.IDGenerator
}

func NewTransactionRepository(db *database.DB, idGenerator id.IDGenerator) TransactionRepository {
	return &trxRepoImpl{
		db:          db,
		idGenerator: idGenerator,
	}
}

func (r *trxRepoImpl) GetTransactionByCampaignID(campaignID uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Model(model.Transaction{}).Where("campaign_id = ?", campaignID).Find(&transactions).Error

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *trxRepoImpl) CreateTransaction(transaction model.Transaction) (model.Transaction, error) {
	id := r.idGenerator.Generate()
	transaction.ID = id
	trx := r.db.Create(&transaction)
	if trx.Error != nil {
		return model.Transaction{}, trx.Error
	}

	return transaction, nil
}
