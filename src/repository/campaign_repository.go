package repository

import (
	"backend-crowdfunding/database"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/sdk/id"
	"context"
)

type CampaignRepository interface {
	FindAllCampaign(c context.Context) ([]model.Campaign, error)
	GetCampaignByUserID(c context.Context, userID string) ([]model.Campaign, error)
	CreateCampaign(c context.Context, campaign model.Campaign) (model.Campaign, error)
	GetCampaignByID(c context.Context, ID string) (model.Campaign, error)
	UpdateCampaign(c context.Context, campaign model.Campaign) (model.Campaign, error)
	UploadImageCampaign(c context.Context, image model.CampaignImage) (model.CampaignImage, error)
	MarkAllImagesAsNonPrimary(c context.Context, campaignID string) (bool, error)
}

type campaignRepoImpl struct {
	db          *database.DB
	idGenerator id.IDGenerator
}

func NewCampaignRepository(db *database.DB, idGenerator id.IDGenerator) CampaignRepository {
	return &campaignRepoImpl{
		db:          db,
		idGenerator: idGenerator,
	}
}

func (r *campaignRepoImpl) FindAllCampaign(c context.Context) ([]model.Campaign, error) {
	var campaigns []model.Campaign
	err := r.db.WithContext(c).Model(model.Campaign{}).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *campaignRepoImpl) GetCampaignByID(c context.Context, ID string) (model.Campaign, error) {
	var campaigns model.Campaign = model.Campaign{ID: ID}
	err := r.db.WithContext(c).Preload("CampaignImages").First(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *campaignRepoImpl) GetCampaignByUserID(c context.Context, userID string) ([]model.Campaign, error) {
	var campaigns []model.Campaign
	err := r.db.WithContext(c).Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = true").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *campaignRepoImpl) CreateCampaign(c context.Context, campaign model.Campaign) (model.Campaign, error) {
	id := r.idGenerator.Generate()
	campaign.ID = id
	data := r.db.WithContext(c).Create(&campaign)
	if data.Error != nil {
		return campaign, data.Error
	}

	return campaign, nil

}

func (r *campaignRepoImpl) UpdateCampaign(c context.Context, campaign model.Campaign) (model.Campaign, error) {
	err := r.db.WithContext(c).Save(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *campaignRepoImpl) UploadImageCampaign(c context.Context, image model.CampaignImage) (model.CampaignImage, error) {
	err := r.db.WithContext(c).Create(&image).Error
	if err != nil {
		return image, err
	}

	return image, nil
}

func (r *campaignRepoImpl) MarkAllImagesAsNonPrimary(c context.Context, campaignID string) (bool, error) {
	err := r.db.WithContext(c).Model(model.CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
