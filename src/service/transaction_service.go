package service

import (
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/repository"
	"backend-crowdfunding/src/request"
	"context"
	"time"
)

type TransactionService interface {
	GetTransactionsByCampaignID(ctx context.Context, campaignID string) ([]model.Transaction, error)
	CreateTransaction(ctx context.Context, input request.CreateTransactionInput) (model.Transaction, error)
}

type trxServiceImpl struct {
	repository repository.TransactionRepository
	timeout    time.Duration
}

func NewTransactionService(repository repository.TransactionRepository) TransactionService {
	return &trxServiceImpl{
		repository: repository,
		timeout:    2 * time.Second,
	}
}

func (s *trxServiceImpl) GetTransactionsByCampaignID(ctx context.Context, campaignID string) ([]model.Transaction, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	data, err := s.repository.GetTransactionByCampaignID(c, campaignID)

	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *trxServiceImpl) CreateTransaction(ctx context.Context, input request.CreateTransactionInput) (model.Transaction, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var trx model.Transaction

	trx.CampaignID = input.CampaignID
	trx.UserID = input.UserID
	trx.Amount = input.Amount
	trx.Code = input.Code
	trx.Status = input.Status

	trx, err := s.repository.CreateTransaction(c, trx)
	if err != nil {
		return trx, err
	}

	return trx, nil
}
