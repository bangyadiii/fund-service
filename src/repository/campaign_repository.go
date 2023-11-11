package repository

import (
	"backend-crowdfunding/database"
	"backend-crowdfunding/insfrastructure/cache"
	ierrors "backend-crowdfunding/sdk/errors"
	"backend-crowdfunding/sdk/id"
	"backend-crowdfunding/src/dto/request"
	"backend-crowdfunding/src/dto/response"
	"backend-crowdfunding/src/model"
	"context"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type CampaignRepository interface {
	FindAllCampaign(c context.Context, params request.CampaignsWithPaginationParam, pgResp response.PaginationResponse) ([]*model.Campaign, *response.PaginationResponse, error)
	GetCampaignByUserID(c context.Context, userID string) ([]*model.Campaign, error)
	CreateCampaign(c context.Context, campaign model.Campaign) (model.Campaign, error)
	GetCampaignByID(c context.Context, ID string) (model.Campaign, error)
	UpdateCampaign(c context.Context, campaign model.Campaign) (model.Campaign, error)
	UploadImageCampaign(c context.Context, image model.CampaignImage) (model.CampaignImage, error)
	MarkAllImagesAsNonPrimary(c context.Context, campaignID string) (bool, error)
}

type campaignRepoImpl struct {
	db          *database.DB
	rd          cache.RedisClient
	idGenerator id.IDGenerator
}

func NewCampaignRepository(db *database.DB, r cache.RedisClient, idGenerator id.IDGenerator) CampaignRepository {
	return &campaignRepoImpl{
		db:          db,
		rd:          r,
		idGenerator: idGenerator,
	}
}

func (r *campaignRepoImpl) FindAllCampaign(c context.Context, params request.CampaignsWithPaginationParam, pgResp response.PaginationResponse) ([]*model.Campaign, *response.PaginationResponse, error) {
	var campaigns []*model.Campaign
	pg := response.FormatPaginationParam(pgResp)
	err := r.db.WithContext(c).
		Model(&campaigns).
		Where(params).
		Count(&pg.TotalElement).Error
	if err != nil {
		return campaigns, nil, err
	}
	if ok := pg.ProcessPagination(); !ok {
		return campaigns, nil, ierrors.NewErrorf(http.StatusNotFound, nil, "campaign not found")
	}

	key := "article:limit:" + strconv.Itoa(int(params.Limit)) + ":page:" + strconv.Itoa(int(params.Page))

	// get data from the cache
	data, err := r.rd.Get(c, key)
	// if there are no data in cache, get data from DB
	if err != nil {
		err = r.db.WithContext(c).
			Model(&campaigns).
			Where(params).
			Preload("CampaignImages", "campaign_images.is_primary = true").
			Offset(int(pg.Offset)).
			Limit(int(pg.Limit)).
			Find(&campaigns).Error

		if err != nil {
			return campaigns, nil, err
		}
		err = r.rd.Set(c, key, campaigns)
		if err != nil {
			return nil, nil, err
		}

		return campaigns, pg, nil
	}
	err = json.Unmarshal([]byte(data), &campaigns)
	if err != nil {
		return campaigns, nil, err
	}
	return campaigns, pg, nil
}

func (r *campaignRepoImpl) GetCampaignByID(c context.Context, ID string) (model.Campaign, error) {
	var campaigns = model.Campaign{ID: ID}
	err := r.db.WithContext(c).
		Preload("CampaignImages").
		First(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *campaignRepoImpl) GetCampaignByUserID(c context.Context, userID string) ([]*model.Campaign, error) {
	var campaigns []*model.Campaign
	err := r.db.WithContext(c).
		Where("user_id = ?", userID).
		Preload("CampaignImages", "campaign_images.is_primary = true").
		Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *campaignRepoImpl) CreateCampaign(c context.Context, campaign model.Campaign) (model.Campaign, error) {
	campaign.ID = r.idGenerator.Generate()

	data := r.db.WithContext(c).Create(&campaign)
	if data.Error != nil {
		return campaign, data.Error
	}

	return campaign, nil

}

func (r *campaignRepoImpl) UpdateCampaign(c context.Context, campaign model.Campaign) (model.Campaign, error) {
	err := r.db.WithContext(c).
		Where("id = ?", campaign.ID).
		Where("updated_at = ?", campaign.UpdatedAt).
		Save(&campaign).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return campaign, ierrors.WrapErrorf(err, 404, nil, "campaign not found")
	} else if err != nil {
		return campaign, ierrors.WrapErrorf(err, 500, nil, err.Error())
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
