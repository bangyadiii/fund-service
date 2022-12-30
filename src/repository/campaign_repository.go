package repository

import (
	"backend-crowdfunding/src/model"

	"gorm.io/gorm"
)

type CampaignRepository interface {
	FindAllCampaign() ([]model.Campaign, error)
	GetCampaignByUserID(userID uint) ([]model.Campaign, error)
	CreateCampaign(campaign model.Campaign) (model.Campaign, error)
	GetCampaignByID(ID uint) (model.Campaign, error)
	UpdateCampaign(campaign model.Campaign) (model.Campaign, error)
	UploadImageCampaign(image model.CampaignImage) (model.CampaignImage, error)
	MarkAllImagesAsNonPrimary(campaignID uint) (bool, error)
}

type campaignRepository struct {
	db *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) *campaignRepository {
	return &campaignRepository{db}
}

func (r *campaignRepository) FindAllCampaign() ([]model.Campaign, error) {
	var campaigns []model.Campaign
	err := r.db.Model(model.Campaign{}).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *campaignRepository) GetCampaignByID(ID uint) (model.Campaign, error) {
	var campaigns model.Campaign = model.Campaign{ID: ID}
	err := r.db.Preload("CampaignImages").First(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *campaignRepository) GetCampaignByUserID(userID uint) ([]model.Campaign, error) {
	var campaigns []model.Campaign
	err := r.db.Where("user_id = ?", userID).Preload("CampainImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *campaignRepository) CreateCampaign(campaign model.Campaign) (model.Campaign, error) {

	data := r.db.Create(&campaign)
	if data.Error != nil {
		return campaign, data.Error
	}

	return campaign, nil

}

func (r *campaignRepository) UpdateCampaign(campaign model.Campaign) (model.Campaign, error) {
	err := r.db.Save(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *campaignRepository) UploadImageCampaign(image model.CampaignImage) (model.CampaignImage, error) {
	err := r.db.Create(&image).Error
	if err != nil {
		return image, err
	}

	return image, nil
}

func (r *campaignRepository) MarkAllImagesAsNonPrimary(campaignID uint) (bool, error) {
	err := r.db.Model(model.CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
