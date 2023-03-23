package handler

import (
	"backend-crowdfunding/helper"
	"backend-crowdfunding/src/request"
	"backend-crowdfunding/src/response"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type CampaignHandler interface {
	GetCampaigns(ctx *gin.Context)
	GetCampaignByID(ctx *gin.Context)
	CreateNewCampaign(ctx *gin.Context)
	UpdateCampaign(ctx *gin.Context)
	UploadCampaignImage(ctx *gin.Context)
}

// GetCampaigns A function that is used to get all campaigns.
func (r *rest) GetCampaigns(ctx *gin.Context) {
	var param request.CampaignsWithPaginationParam
	err := ctx.ShouldBindQuery(&param)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
		return
	}

	data, pg, err := r.service.Campaign.GetCampaigns(ctx.Request.Context(), param)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
		return
	}

	response.SuccessResponseWithPagination(ctx, http.StatusOK, "OK", data, pg)
}

// GetCampaignByID A function that is used to get a campaign by ID.
func (r *rest) GetCampaignByID(ctx *gin.Context) {
	var input request.GetCampaignByIDInput
	err := ctx.ShouldBindUri(&input)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
		return
	}

	data, err := r.service.Campaign.GetCampaignByID(ctx.Request.Context(), input)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
		return
	}
	response.SuccessResponse(ctx, http.StatusOK, "OK", data)
}

// CreateNewCampaign A function that is used to create a new campaign.
func (r *rest) CreateNewCampaign(ctx *gin.Context) {
	var input request.CreateCampaignInput
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		errorData := helper.FormatErrorValidation(err)
		response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", errorData)
		return
	}
	userString, ok := ctx.Get("current_user")
	if !ok {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Unautheticated", nil)
	}
	userResp, ok := userString.(response.UserResponse)
	if !ok {
		response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", nil)
	}

	input.User.ID = userResp.ID
	input.User.Email = userResp.Email

	campaignRes, err := r.service.Campaign.CreateCampaign(ctx.Request.Context(), input)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
		return
	}
	response.SuccessResponse(ctx, http.StatusCreated, "CREATED", campaignRes)
}

// UpdateCampaign A function that is used to update a campaign.
func (r *rest) UpdateCampaign(ctx *gin.Context) {
	var input request.UpdateCampaignInput
	var campaignID request.GetCampaignByIDInput

	err := ctx.ShouldBindUri(&campaignID)

	if err != nil {
		errorData := helper.FormatErrorValidation(err)
		response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", errorData)
		return
	}

	err = ctx.ShouldBindJSON(&input)

	if err != nil {
		errorData := helper.FormatErrorValidation(err)
		response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", errorData)
		return
	}
	userString, ok := ctx.Get("current_user")
	if !ok {
		response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", nil)
	}
	userResp, ok := userString.(response.UserResponse)
	if !ok {
		response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", nil)
	}

	input.User.ID = userResp.ID
	input.User.Email = userResp.Email

	data, err := r.service.Campaign.UpdateCampaign(ctx.Request.Context(), campaignID, input)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "OK", data)
}

func (r *rest) UploadCampaignImage(ctx *gin.Context) {
	var input request.UploadCampaignImageInput
	err := ctx.ShouldBind(&input)

	if err != nil {
		err := helper.FormatErrorValidation(err)
		res := response.APIResponse("Bad Reqeust", http.StatusBadRequest, "error", nil, err)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	userResp, ok := ctx.MustGet("current_user").(response.UserResponse)
	if !ok {
		response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", nil)
	}

	input.User.ID = userResp.ID
	input.User.Email = userResp.Email

	imageFile, err := ctx.FormFile("campaign_image")

	if err != nil {
		res := response.APIResponse("Bad Reqeust", http.StatusBadRequest, "error", nil, err)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	path := fmt.Sprintf("assets/images/campaigns/%s-%d-%s", input.CampaignID, time.Now().Day(), imageFile.Filename)

	err = ctx.SaveUploadedFile(imageFile, path)

	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		apiResponse := response.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data, err.Error())
		ctx.JSON(http.StatusBadRequest, apiResponse)
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
		res := response.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data, err.Error())
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "OK", payload)
}
