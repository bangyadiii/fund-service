package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Campaign, error)
	GetCampaignByUserID( userID int) ([]Campaign, error)
}

type repository struct {
	db *gorm.DB
}


func NewRepository(db *gorm.DB) *repository{
	return &repository{db}
}

func (r *repository) FindAll() ([] Campaign, error){
	var campaigns []Campaign
	err := r.db.Preload("CampaignImages", "campaign_primary.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}


func (r *repository) GetCampaignByUserID(userID int) ([] Campaign, error){
	var campaigns []Campaign
	err := r.db.Where("user_id = ?", userID).Preload("CampainImages").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}