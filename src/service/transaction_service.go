package service

import (
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/repository"
	"backend-crowdfunding/src/request"
	"backend-crowdfunding/src/response"
	"context"
	"time"
)

type TransactionService interface {
	GetTransactionsByCampaignID(ctx context.Context, campaignID string) ([]response.TransactionResponse, error)
	CreateTransaction(ctx context.Context, input request.CreateTransactionInput) (response.TransactionResponse, error)
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

func (s *trxServiceImpl) GetTransactionsByCampaignID(ctx context.Context, campaignID string) ([]response.TransactionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	var trxResponse []response.TransactionResponse

	data, err := s.repository.GetTransactionByCampaignID(c, campaignID)

	if err != nil {
		return trxResponse, err
	}

	trxResponse = response.FormatTransactionCollections(data)

	return trxResponse, nil
}

// CreateTransaction A function that will be called by the controller.
func (s *trxServiceImpl) CreateTransaction(ctx context.Context, input request.CreateTransactionInput) (response.TransactionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	var trx model.Transaction
	var trxResponse response.TransactionResponse

	trx.CampaignID = input.CampaignID
	trx.UserID = input.UserID
	trx.Amount = input.Amount
	trx.Code = input.Code
	trx.Status = input.Status

	trx, err := s.repository.CreateTransaction(c, trx)
	if err != nil {
		return trxResponse, err
	}
	trxResponse = response.FormatTransaction(trx)
	return trxResponse, nil
}
