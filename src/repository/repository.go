package repository

import (
	"backend-crowdfunding/database"
	"backend-crowdfunding/src/util/id"
)

type Repository struct {
	CampaignRepo CampaignRepository
	TrxRepo      TransactionRepository
	UserRepo     UserRepository
}

func InitRepository(db *database.DB, idGenerator id.IDGenerator) *Repository {
	return &Repository{
		CampaignRepo: NewCampaignRepository(db, idGenerator),
		TrxRepo:      NewTransactionRepository(db, idGenerator),
		UserRepo:     NewUserRepository(db, idGenerator),
	}
}
