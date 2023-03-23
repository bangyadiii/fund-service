package repository

import (
	"backend-crowdfunding/database"
	"backend-crowdfunding/insfrastructure/cache"
	"backend-crowdfunding/insfrastructure/firebase"
	"backend-crowdfunding/sdk/id"
)

type Repository struct {
	CampaignRepo CampaignRepository
	TrxRepo      TransactionRepository
	UserRepo     UserRepository
}

func InitRepository(db *database.DB, r cache.RedisClient, idGenerator id.IDGenerator, firebase *firebase.Firebase) *Repository {
	return &Repository{
		CampaignRepo: NewCampaignRepository(db, r, idGenerator),
		TrxRepo:      NewTransactionRepository(db, idGenerator),
		UserRepo:     NewUserRepository(db, firebase, idGenerator),
	}
}
