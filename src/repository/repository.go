package repository

import "backend-crowdfunding/database"

type Repository struct {
	CampaignRepo CampaignRepository
	TrxRepo      TransactionRepository
	UserRepo     UserRepository
}

func InitRepository(db *database.DB) *Repository {
	return &Repository{
		CampaignRepo: NewCampaignRepository(db),
		TrxRepo:      NewTransactionRepository(db),
		UserRepo:     NewUserRepository(db),
	}
}
