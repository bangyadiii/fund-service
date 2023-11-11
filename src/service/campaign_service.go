package service

import (
	ierrors "backend-crowdfunding/sdk/errors"
	"backend-crowdfunding/src/dto/request"
	"backend-crowdfunding/src/dto/response"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/repository"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gosimple/slug"
)

type CampaignService interface {
	GetCampaigns(c context.Context, params request.CampaignsWithPaginationParam) ([]response.CampaignResponse, *response.PaginationResponse, error)
	CreateCampaign(c context.Context, input request.CreateCampaignInput, userID string) (response.CampaignResponse, error)
	GetCampaignByID(c context.Context, input request.GetCampaignByIDInput) (response.CampaignDetailFormatter, error)
	UpdateCampaign(c context.Context, campaignID request.GetCampaignByIDInput, input request.UpdateCampaignInput) (response.CampaignResponse, error)
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

func (s *campaignServiceImpl) GetCampaigns(c context.Context, params request.CampaignsWithPaginationParam) ([]response.CampaignResponse, *response.PaginationResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	var campaignResponse []response.CampaignResponse
	var pg = response.ConvertPaginationParamToPaginationResponse(params.PaginationParam)

	if params.UserID != "" {
		campaigns, err := s.repository.GetCampaignByUserID(ctx, params.UserID)
		if err != nil {
			return campaignResponse, pg, err
		}
		return response.FormatCampaignCollections(campaigns), pg, nil
	}

	campaigns, pg, err := s.repository.FindAllCampaign(ctx, params, *pg)
	if err != nil {
		return campaignResponse, pg, err
	}

	return response.FormatCampaignCollections(campaigns), pg, nil

}

func (s *campaignServiceImpl) GetCampaignByID(c context.Context, input request.GetCampaignByIDInput) (response.CampaignDetailFormatter, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	var campaignRes response.CampaignDetailFormatter

	campaign, err := s.repository.GetCampaignByID(ctx, input.ID)

	if err != nil {
		return campaignRes, err
	}

	campaignRes = response.FormatCampaignDetail(campaign)
	return campaignRes, nil
}

func (s *campaignServiceImpl) CreateCampaign(c context.Context, input request.CreateCampaignInput, userID string) (response.CampaignResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	var campaignRes response.CampaignResponse

	campaign := model.Campaign{}

	campaign.Name = input.Name
	campaign.UserID = userID
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.BackerCount = int(input.BackerCount)
	campaign.GoalAmount = int(input.GoalAmount)

	slugCandidate := fmt.Sprintf("%s %d", input.Name, rand.Int())
	campaign.Slug = slug.Make(slugCandidate)

	newCampaign, err := s.repository.CreateCampaign(ctx, campaign)

	if err != nil {
		return campaignRes, err
	}

	return response.FormatCampaign(newCampaign), nil
}

func (s *campaignServiceImpl) UpdateCampaign(c context.Context, campaignID request.GetCampaignByIDInput, input request.UpdateCampaignInput) (response.CampaignResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	var campaignRes response.CampaignResponse

	oldCampaign, err := s.repository.GetCampaignByID(ctx, campaignID.ID)
	if err != nil {
		return campaignRes, err
	}

	if oldCampaign.UserID != input.User.ID {
		return campaignRes, ierrors.NewErrorf(403, nil, "unauthorized")
	}

	oldCampaign.Description = input.Description
	oldCampaign.ShortDescription = input.ShortDescription
	oldCampaign.BackerCount = int(input.BackerCount)
	oldCampaign.CurrentAmount = int(input.CurrentAmount)
	oldCampaign.GoalAmount = int(input.GoalAmount)
	oldCampaign.Perks = input.Perks

	data, err := s.repository.UpdateCampaign(ctx, oldCampaign)

	if err != nil {
		return campaignRes, err
	}

	campaignRes = response.FormatCampaign(data)
	return campaignRes, nil

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
		return image, ierrors.NewErrorf(403, nil, "unauthorized")
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
