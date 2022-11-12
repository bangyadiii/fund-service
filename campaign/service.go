package campaign

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID uint) ([]Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	GetCampaignByID(input GetCampaignByIDInput) (Campaign, error)
	UpdateCampaign(campaignID GetCampaignByIDInput, input UpdateCampaignInput) (Campaign, error)
	UploadCampaignImage(input UploadCampaignImageInput) (CampaignImage, error)
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

func (s *service) UpdateCampaign(campaignID GetCampaignByIDInput, input UpdateCampaignInput) (Campaign, error) {
	oldCampaign, err := s.repository.GetCampaignByID(campaignID.ID)
	if err != nil {
		return oldCampaign, err
	}

	if oldCampaign.UserID != input.User.ID {
		return oldCampaign, errors.New("Unauthorized")
	}

	oldCampaign.Description = input.Description
	oldCampaign.ShortDescription = input.ShortDescription
	oldCampaign.BackerCount = int(input.BackerCount)
	oldCampaign.CurrentAmount = int(input.CurrentAmount)
	oldCampaign.GoalAmount = int(input.GoalAmount)
	oldCampaign.Perks = input.Perks

	data, err := s.repository.Update(oldCampaign)

	if err != nil {
		return data, err
	}

	return data, nil

}

func (s *service) UploadCampaignImage(input UploadCampaignImageInput) (CampaignImage, error) {
	var image CampaignImage
	image.CampaignID = input.CampaignID
	image.ImageName = input.ImageName
	image.IsPrimary = input.IsPrimary

	campaignImage, err := s.repository.UploadImage(image)
	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
}
