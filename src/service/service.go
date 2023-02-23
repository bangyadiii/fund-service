package service

import (
	"backend-crowdfunding/config"
	"backend-crowdfunding/src/repository"
)

type Service struct {
	Auth     AuthService
	User     UserService
	Campaign CampaignService
	Trx      TransactionService
}

func InitService(configuration config.Config, repo *repository.Repository) *Service {
	return &Service{
		Auth:     NewAuthService(),
		User:     NewUserService(repo.UserRepo),
		Campaign: NewCampaignService(repo.CampaignRepo),
		Trx:      NewTransactionService(repo.TrxRepo),
	}
}
