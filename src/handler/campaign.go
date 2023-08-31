package handler

import (
	"backend-crowdfunding/helper"
	"backend-crowdfunding/src/request"
	"backend-crowdfunding/src/response"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
	"time"
)

type CampaignHandler interface {
	GetCampaigns(ctx *fiber.Ctx) error
	GetCampaignByID(ctx *fiber.Ctx) error
	CreateNewCampaign(ctx *fiber.Ctx) error
	UpdateCampaign(ctx *fiber.Ctx) error
	UploadCampaignImage(ctx *fiber.Ctx) error
}

// GetCampaigns A function that is used to get all campaigns.
func (r *Rest) GetCampaigns(ctx *fiber.Ctx) error {
	var param request.CampaignsWithPaginationParam
	err := ctx.ParamsParser(&param)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	data, pg, err := r.service.Campaign.GetCampaigns(ctx.Context(), param)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return response.SuccessResponseWithPagination(ctx, http.StatusOK, "OK", data, pg)
}

// GetCampaignByID A function that is used to get a campaign by ID.
func (r *Rest) GetCampaignByID(ctx *fiber.Ctx) error {
	var input request.GetCampaignByIDInput
	err := ctx.ParamsParser(&input)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	data, err := r.service.Campaign.GetCampaignByID(ctx.Context(), input)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}
	return response.SuccessResponse(ctx, http.StatusOK, "OK", data)
}

// CreateNewCampaign A function that is used to create a new campaign.
func (r *Rest) CreateNewCampaign(ctx *fiber.Ctx) error {
	var input request.CreateCampaignInput
	err := ctx.BodyParser(&input)
	if err != nil {
		errorData := helper.FormatErrorValidation(err)
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", errorData)
	}
	userString := ctx.Locals("current_user")
	userResp, ok := userString.(response.UserResponse)
	if !ok {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", nil)
	}

	input.User.ID = userResp.ID
	input.User.Email = userResp.Email

	campaignRes, err := r.service.Campaign.CreateCampaign(ctx.Context(), input)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}
	return response.SuccessResponse(ctx, http.StatusCreated, "CREATED", campaignRes)
}

// UpdateCampaign A function that is used to update a campaign.
func (r *Rest) UpdateCampaign(ctx *fiber.Ctx) error {
	var input request.UpdateCampaignInput
	var campaignID request.GetCampaignByIDInput

	err := ctx.BodyParser(&campaignID)

	if err != nil {
		errorData := helper.FormatErrorValidation(err)
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", errorData)
	}

	err = ctx.BodyParser(&input)

	if err != nil {
		errorData := helper.FormatErrorValidation(err)
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", errorData)
	}
	userString := ctx.Locals("current_user")

	userResp, ok := userString.(response.UserResponse)

	if !ok {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", nil)
	}

	input.User.ID = userResp.ID
	input.User.Email = userResp.Email

	data, err := r.service.Campaign.UpdateCampaign(ctx.Context(), campaignID, input)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error", err.Error())
	}

	return response.SuccessResponse(ctx, http.StatusOK, "OK", data)
}

func (r *Rest) UploadCampaignImage(ctx *fiber.Ctx) error {
	var input request.UploadCampaignImageInput
	err := ctx.BodyParser(&input)

	if err != nil {
		err := helper.FormatErrorValidation(err)
		res := response.APIResponse("Bad Reqeust", http.StatusBadRequest, "error", nil, err)
		return ctx.Status(http.StatusBadRequest).JSON(res)
	}
	userResp, ok := ctx.Locals("current_user").(response.UserResponse)
	if !ok {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", nil)
	}

	input.User.ID = userResp.ID
	input.User.Email = userResp.Email

	imageFile, err := ctx.FormFile("campaign_image")

	if err != nil {
		res := response.APIResponse("Bad Reqeust", http.StatusBadRequest, "error", nil, err)
		return ctx.Status(http.StatusBadRequest).JSON(res)
	}

	path := fmt.Sprintf("assets/images/campaigns/%s-%d-%s", input.CampaignID, time.Now().Day(), imageFile.Filename)

	err = ctx.SaveFile(imageFile, path)

	if err != nil {
		data := fiber.Map{
			"is_uploaded": false,
		}
		apiResponse := response.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data, err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(apiResponse)
	}

	input.ImageName = path
	payload, err := r.service.Campaign.UploadCampaignImage(ctx.Context(), input)

	if err != nil {
		err := os.Remove(path)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(err)
		}
		data := fiber.Map{
			"is_uploaded": false,
		}
		res := response.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data, err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(res)
	}

	return response.SuccessResponse(ctx, http.StatusOK, "OK", payload)
}
