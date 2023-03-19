package handler

import (
	"backend-crowdfunding/helper"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/request"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type CampaignHandler interface {
	GetCampaigns(c *gin.Context)
	GetCampaignByID(c *gin.Context)
	CreateNewCampaign(c *gin.Context)
	UpdateCampaign(c *gin.Context)
	UploadCampaignImage(c *gin.Context)
}

// GetCampaigns A function that is used to get all campaigns.
func (r *rest) GetCampaigns(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	data, err := r.service.Campaign.GetCampaigns(ctx.Request.Context(), userID)

	if err != nil {
		helper.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
		return
	}
	helper.SuccessResponse(ctx, http.StatusOK, "OK", data)
}

// GetCampaignByID A function that is used to get a campaign by ID.
func (r *rest) GetCampaignByID(ctx *gin.Context) {
	var input request.GetCampaignByIDInput
	err := ctx.ShouldBindUri(&input)

	if err != nil {
		helper.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
		return
	}

	data, err := r.service.Campaign.GetCampaignByID(ctx.Request.Context(), input)

	if err != nil {
		helper.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
		return
	}
	helper.SuccessResponse(ctx, http.StatusOK, "OK", data)
}

// CreateNewCampaign A function that is used to create a new campaign.
func (r *rest) CreateNewCampaign(c *gin.Context) {
	var input request.CreateCampaignInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorData := helper.FormatErrorValidation(err)
		helper.ErrorResponse(c, http.StatusBadRequest, "BAD REQUEST", errorData)
		return
	}
	curUser := c.MustGet("current_user").(model.User)
	input.User = curUser

	campaignRes, err := r.service.Campaign.CreateCampaign(c.Request.Context(), input)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "BAD REQUEST", err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "CREATED", campaignRes)
}

// UpdateCampaign A function that is used to update a campaign.
func (r *rest) UpdateCampaign(ctx *gin.Context) {
	var input request.UpdateCampaignInput
	var campaignID request.GetCampaignByIDInput

	err := ctx.ShouldBindUri(&campaignID)

	if err != nil {
		errorData := helper.FormatErrorValidation(err)
		helper.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", errorData)
		return
	}

	err = ctx.ShouldBindJSON(&input)

	if err != nil {
		errorData := helper.FormatErrorValidation(err)
		helper.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", errorData)
		return
	}
	input.User = ctx.MustGet("current_user").(model.User)

	data, err := r.service.Campaign.UpdateCampaign(ctx.Request.Context(), campaignID, input)

	if err != nil {
		helper.ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	helper.SuccessResponse(ctx, http.StatusOK, "OK", data)
}

func (r *rest) UploadCampaignImage(ctx *gin.Context) {
	var input request.UploadCampaignImageInput
	err := ctx.ShouldBind(&input)

	if err != nil {
		err := helper.FormatErrorValidation(err)
		res := helper.APIResponse("Bad Reqeust", http.StatusBadRequest, "error", nil, err)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	input.User = ctx.MustGet("current_user").(model.User)

	imageFile, err := ctx.FormFile("campaign_image")

	if err != nil {
		res := helper.APIResponse("Bad Reqeust", http.StatusBadRequest, "error", nil, err)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	path := fmt.Sprintf("assets/images/campaigns/%s-%d-%s", input.CampaignID, time.Now().Day(), imageFile.Filename)

	err = ctx.SaveUploadedFile(imageFile, path)

	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data, err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	input.ImageName = path
	payload, err := r.service.Campaign.UploadCampaignImage(ctx.Request.Context(), input)

	if err != nil {
		err := os.Remove(path)
		if err != nil {
			ctx.JSON(500, err)
		}
		data := gin.H{
			"is_uploaded": false,
		}
		res := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data, err.Error())
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	helper.SuccessResponse(ctx, http.StatusOK, "OK", payload)
}
