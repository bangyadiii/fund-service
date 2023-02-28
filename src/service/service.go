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
	authService := NewAuthService(&configuration)
	return &Service{
		Auth:     authService,
		User:     NewUserService(repo.UserRepo, authService),
		Campaign: NewCampaignService(repo.CampaignRepo),
		Trx:      NewTransactionService(repo.TrxRepo),
	}
}
