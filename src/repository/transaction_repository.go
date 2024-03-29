package repository

import (
	"backend-crowdfunding/database"
	"backend-crowdfunding/sdk/id"
	model "backend-crowdfunding/src/model"
	"context"
)

type TransactionRepository interface {
	GetTransactionByCampaignID(ctx context.Context, campaignID string) ([]model.Transaction, error)
	CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error)
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

func (r *trxRepoImpl) GetTransactionByCampaignID(ctx context.Context, campaignID string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.WithContext(ctx).Model(model.Transaction{}).Where("campaign_id = ?", campaignID).Find(&transactions).Error

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *trxRepoImpl) CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	transaction.ID = r.idGenerator.Generate()
	trx := r.db.WithContext(ctx).Create(&transaction)
	if trx.Error != nil {
		return model.Transaction{}, trx.Error
	}

	return transaction, nil
}
