package transaction

type Service interface {
	GetTransactionsByCampaignID(campaignID uint) ([]Transaction, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) GetTransactionsByCampaignID(campaignID uint) ([]Transaction, error) {

	data, err := s.repository.GetTransactionByCampaignID(campaignID)

	if err != nil {
		return data, err
	}

	return data, nil
}
