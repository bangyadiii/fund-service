package service

import (
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/repository"
	"backend-crowdfunding/src/request"
)

type TransactionService interface {
	GetTransactionsByCampaignID(campaignID uint) ([]model.Transaction, error)
	CreateTransaction(input request.CreateTransactionInput) (model.Transaction, error)
}

type trxServiceImpl struct {
	repository repository.TransactionRepository
}

func NewTransactionService(repository repository.TransactionRepository) TransactionService {
	return &trxServiceImpl{repository: repository}
}

func (s *trxServiceImpl) GetTransactionsByCampaignID(campaignID uint) ([]model.Transaction, error) {

	data, err := s.repository.GetTransactionByCampaignID(campaignID)

	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *trxServiceImpl) CreateTransaction(input request.CreateTransactionInput) (model.Transaction, error) {
	var trx model.Transaction

	trx.CampaignID = input.CampaignID
	trx.UserID = input.UserID
	trx.Amount = input.Amount
	trx.Code = input.Code
	trx.Status = input.Status

	trx, err := s.repository.CreateTransaction(trx)
	if err != nil {
		return trx, err
	}

	return trx, nil
}
