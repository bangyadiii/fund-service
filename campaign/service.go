package campaign

import (
	"fmt"
	"math/rand"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID uint) ([]Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	GetCampaignByID(input GetCampaignByIDInput) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userID uint) ([]Campaign, error) {
	if userID != 0 {
		campaigns, err := s.repository.GetCampaignByUserID(userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}
	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil

}

func (s *service) GetCampaignByID(input GetCampaignByIDInput) (Campaign, error) {
	campaign, err := s.repository.GetCampaignByID(input.ID)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}

	campaign.Name = input.Name
	campaign.UserID = input.User.ID
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.BackerCount = int(input.BackerCount)
	campaign.GoalAmount = int(input.GoalAmount)

	slugCandidate := fmt.Sprintf("%s %d%d", input.Name, input.User.ID, rand.Int())
	campaign.Slug = slug.Make(slugCandidate)

	newCampaign, err := s.repository.Create(campaign)

	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil

}
