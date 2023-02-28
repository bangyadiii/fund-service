package repository

import (
	"backend-crowdfunding/database"
	"backend-crowdfunding/insfrastructure/firebase"
	"backend-crowdfunding/sdk/id"
)

type Repository struct {
	CampaignRepo CampaignRepository
	TrxRepo      TransactionRepository
	UserRepo     UserRepository
}

func InitRepository(db *database.DB, idGenerator id.IDGenerator, firebase *firebase.Firebase) *Repository {
	return &Repository{
		CampaignRepo: NewCampaignRepository(db, idGenerator),
		TrxRepo:      NewTransactionRepository(db, idGenerator),
		UserRepo:     NewUserRepository(db, firebase, idGenerator),
	}
}
