package service

import (
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/repository"
	"backend-crowdfunding/src/request"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gosimple/slug"
)

type CampaignService interface {
	GetCampaigns(c context.Context, userID string) ([]model.Campaign, error)
	CreateCampaign(c context.Context, input request.CreateCampaignInput) (model.Campaign, error)
	GetCampaignByID(c context.Context, input request.GetCampaignByIDInput) (model.Campaign, error)
	UpdateCampaign(c context.Context, campaignID request.GetCampaignByIDInput, input request.UpdateCampaignInput) (model.Campaign, error)
	UploadCampaignImage(c context.Context, input request.UploadCampaignImageInput) (model.CampaignImage, error)
}

type campaignServiceImpl struct {
	repository repository.CampaignRepository
	timeout    time.Duration
}

func NewCampaignService(repository repository.CampaignRepository) CampaignService {
	return &campaignServiceImpl{
		repository: repository,
		timeout:    2 * time.Second,
	}
}

func (s *campaignServiceImpl) GetCampaigns(c context.Context, userID string) ([]model.Campaign, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	if userID != "" {
		campaigns, err := s.repository.GetCampaignByUserID(ctx, userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}
	campaigns, err := s.repository.FindAllCampaign(ctx)
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil

}

func (s *campaignServiceImpl) GetCampaignByID(c context.Context, input request.GetCampaignByIDInput) (model.Campaign, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	campaign, err := s.repository.GetCampaignByID(ctx, input.ID)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *campaignServiceImpl) CreateCampaign(c context.Context, input request.CreateCampaignInput) (model.Campaign, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	campaign := model.Campaign{}

	campaign.Name = input.Name
	campaign.UserID = input.User.ID
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.BackerCount = int(input.BackerCount)
	campaign.GoalAmount = int(input.GoalAmount)

	slugCandidate := fmt.Sprintf("%s %d", input.Name, rand.Int())
	campaign.Slug = slug.Make(slugCandidate)

	newCampaign, err := s.repository.CreateCampaign(ctx, campaign)

	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s *campaignServiceImpl) UpdateCampaign(c context.Context, campaignID request.GetCampaignByIDInput, input request.UpdateCampaignInput) (model.Campaign, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	oldCampaign, err := s.repository.GetCampaignByID(ctx, campaignID.ID)
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

	data, err := s.repository.UpdateCampaign(ctx, oldCampaign)

	if err != nil {
		return data, err
	}

	return data, nil

}

func (s *campaignServiceImpl) UploadCampaignImage(c context.Context, input request.UploadCampaignImageInput) (model.CampaignImage, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	var image model.CampaignImage
	campaign, err := s.repository.GetCampaignByID(ctx, input.CampaignID)
	if err != nil {
		return image, err
	}

	if campaign.User.ID != input.User.ID {
		return image, errors.New("Unauthorized")
	}

	if input.IsPrimary {
		_, err := s.repository.MarkAllImagesAsNonPrimary(ctx, input.CampaignID)

		if err != nil {
			return image, err
		}

	}

	image.IsPrimary = input.IsPrimary
	image.CampaignID = input.CampaignID
	image.ImageName = input.ImageName

	campaignImage, err := s.repository.UploadImageCampaign(ctx, image)
	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
}
