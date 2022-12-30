package service

import (
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/repository"
	"backend-crowdfunding/src/request"
	"errors"
	"fmt"
	"math/rand"

	"github.com/gosimple/slug"
)

type CampaignService interface {
	GetCampaigns(userID uint) ([]model.Campaign, error)
	CreateCampaign(input request.CreateCampaignInput) (model.Campaign, error)
	GetCampaignByID(input request.GetCampaignByIDInput) (model.Campaign, error)
	UpdateCampaign(campaignID request.GetCampaignByIDInput, input request.UpdateCampaignInput) (model.Campaign, error)
	UploadCampaignImage(input request.UploadCampaignImageInput) (model.CampaignImage, error)
}

type campaignService struct {
	repository repository.CampaignRepository
}

func NewCampaignService(repository repository.CampaignRepository) *campaignService {
	return &campaignService{repository}
}

func (s *campaignService) GetCampaigns(userID uint) ([]model.Campaign, error) {
	if userID != 0 {
		campaigns, err := s.repository.GetCampaignByUserID(userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}
	campaigns, err := s.repository.FindAllCampaign()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil

}

func (s *campaignService) GetCampaignByID(input request.GetCampaignByIDInput) (model.Campaign, error) {
	campaign, err := s.repository.GetCampaignByID(input.ID)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *campaignService) CreateCampaign(input request.CreateCampaignInput) (model.Campaign, error) {
	campaign := model.Campaign{}

	campaign.Name = input.Name
	campaign.UserID = input.User.ID
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.BackerCount = int(input.BackerCount)
	campaign.GoalAmount = int(input.GoalAmount)

	slugCandidate := fmt.Sprintf("%s %d%d", input.Name, input.User.ID, rand.Int())
	campaign.Slug = slug.Make(slugCandidate)

	newCampaign, err := s.repository.CreateCampaign(campaign)

	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s *campaignService) UpdateCampaign(campaignID request.GetCampaignByIDInput, input request.UpdateCampaignInput) (model.Campaign, error) {
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

	data, err := s.repository.UpdateCampaign(oldCampaign)

	if err != nil {
		return data, err
	}

	return data, nil

}

func (s *campaignService) UploadCampaignImage(input request.UploadCampaignImageInput) (model.CampaignImage, error) {
	var image model.CampaignImage
	campaign, err := s.repository.GetCampaignByID(input.CampaignID)
	if err != nil {
		return image, err
	}

	if campaign.User.ID != input.User.ID {
		return image, errors.New("Unauthorized")
	}

	if input.IsPrimary {
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)

		if err != nil {
			return image, err
		}

	}

	image.IsPrimary = input.IsPrimary
	image.CampaignID = input.CampaignID
	image.ImageName = input.ImageName

	campaignImage, err := s.repository.UploadImageCampaign(image)
	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
}
