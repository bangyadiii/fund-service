package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Campaign, error)
	GetCampaignByUserID(userID int) ([]Campaign, error)
	Create(campaign Campaign) (Campaign, error)
	GetCampaignByID(ID uint32) (Campaign, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Model(Campaign{}).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) GetCampaignByID(ID uint32) (Campaign, error) {
	var campaigns Campaign
	err := r.db.Preload("CampainImages").First(&campaigns, ID).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) GetCampaignByUserID(userID int) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Where("user_id = ?", userID).Preload("CampainImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) Create(campaign Campaign) (Campaign, error) {

	data := r.db.Create(&campaign)
	if data.Error != nil {
		return campaign, data.Error
	}

	return campaign, nil

}
