package repository

import (
	"backend-crowdfunding/database"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/util/id"
)

type CampaignRepository interface {
	FindAllCampaign() ([]model.Campaign, error)
	GetCampaignByUserID(userID uint) ([]model.Campaign, error)
	CreateCampaign(campaign model.Campaign) (model.Campaign, error)
	GetCampaignByID(ID string) (model.Campaign, error)
	UpdateCampaign(campaign model.Campaign) (model.Campaign, error)
	UploadImageCampaign(image model.CampaignImage) (model.CampaignImage, error)
	MarkAllImagesAsNonPrimary(campaignID string) (bool, error)
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

func (r *campaignRepoImpl) FindAllCampaign() ([]model.Campaign, error) {
	var campaigns []model.Campaign
	err := r.db.Model(model.Campaign{}).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *campaignRepoImpl) GetCampaignByID(ID string) (model.Campaign, error) {
	var campaigns model.Campaign = model.Campaign{ID: ID}
	err := r.db.Preload("CampaignImages").First(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *campaignRepoImpl) GetCampaignByUserID(userID uint) ([]model.Campaign, error) {
	var campaigns []model.Campaign
	err := r.db.Where("user_id = ?", userID).Preload("CampainImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *campaignRepoImpl) CreateCampaign(campaign model.Campaign) (model.Campaign, error) {
	id := r.idGenerator.Generate()
	campaign.ID = id
	data := r.db.Create(&campaign)
	if data.Error != nil {
		return campaign, data.Error
	}

	return campaign, nil

}

func (r *campaignRepoImpl) UpdateCampaign(campaign model.Campaign) (model.Campaign, error) {
	err := r.db.Save(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *campaignRepoImpl) UploadImageCampaign(image model.CampaignImage) (model.CampaignImage, error) {
	err := r.db.Create(&image).Error
	if err != nil {
		return image, err
	}

	return image, nil
}

func (r *campaignRepoImpl) MarkAllImagesAsNonPrimary(campaignID string) (bool, error) {
	err := r.db.Model(model.CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
