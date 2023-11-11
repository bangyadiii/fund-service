package handler

import (
	"backend-crowdfunding/sdk/errors"
	"backend-crowdfunding/src/dto/request"
	"backend-crowdfunding/src/dto/response"
	"backend-crowdfunding/src/validator"
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
// @Summary Get all campaigns.
// @Tags Campaign
// @Accept application/json
// @Produce json
// @Param q query request.PaginationParam false "query params"
// @Router /campaigns/ [get]
// @Success 200 {object} response.WithPagination{meta=response.Meta,data=[]response.CampaignResponse} "Success"
// @Failure 400 {object} response.Response{meta=response.Meta,errors=request.CampaignsWithPaginationParam} "Validation Error"
// @Failure 500 {object} response.Response{meta=response.Meta} "Internal Server Error"
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
// @Summary Get Campaign By ID.
// @Tags Campaign
// @Accept application/json
// @Produce json
// @Param id path string true "Campaign ID"
// @Router /campaigns/{campaign_id} [get]
// @Success 200 {object} response.WithPagination{meta=response.Meta,data=response.CampaignResponse} "Success"
// @Failure 400 {object} response.Response{meta=response.Meta,errors=request.GetCampaignByIDInput} "Validation Error"
// @Failure 500 {object} response.Response{meta=response.Meta} "Internal Server Error"
func (r *Rest) GetCampaignByID(ctx *fiber.Ctx) error {
	var input request.GetCampaignByIDInput
	err := ctx.ParamsParser(&input)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	if errMap := validator.Validate(input); errMap != nil {
		err := errors.NewErrorf(http.StatusBadRequest, errMap, "validation error")
		return response.RenderErrorResponse(ctx, "validation error", err)
	}

	data, err := r.service.Campaign.GetCampaignByID(ctx.Context(), input)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return response.SuccessResponse(ctx, http.StatusOK, "OK", data)
}

// CreateNewCampaign A function that is used to create a new campaign.
// @Summary Create Campaign.
// @Tags Campaign
// @Accept application/json
// @Produce json
// @Param request body request.CreateCampaignInput true "Campaign request"
// @Param Authorization header string true "Access token"
// @Router /campaigns/ [post]
// @Success 201 {object} response.Response{meta=response.Meta,data=response.CampaignResponse} "Success"
// @Failure 400 {object} response.Response{meta=response.Meta,errors=request.CreateCampaignInput} "Validation Error"
// @Failure 500 {object} response.Response{meta=response.Meta} "Internal Server Error"
func (r *Rest) CreateNewCampaign(ctx *fiber.Ctx) error {
	var input request.CreateCampaignInput
	err := ctx.BodyParser(&input)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	if errMap := validator.Validate(input); errMap != nil {
		err := errors.NewErrorf(http.StatusBadRequest, errMap, "validation error")
		return response.RenderErrorResponse(ctx, "validation error", err)
	}

	userString := ctx.Locals("current_user")
	userResp, ok := userString.(response.UserResponse)
	if !ok {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", nil)
	}

	campaignRes, err := r.service.Campaign.CreateCampaign(ctx.Context(), input, userResp.ID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}
	return response.SuccessResponse(ctx, http.StatusCreated, "CREATED", campaignRes)
}

// UpdateCampaign A function that is used to update a campaign.
// @Summary Update Campaign.
// @Tags Campaign
// @Accept application/json
// @Produce json
// @Param Authorization header string true "Access token"
// @Param id path string true "Campaign ID"
// @Param request body request.UpdateCampaignInput true "Campaign payload"
// @Router /campaigns/{campaign_id} [put]
// @Success 200 {object} response.Response{meta=response.Meta,data=response.CampaignResponse} "Success"
// @Failure 400 {object} response.Response{meta=response.Meta,errors=request.CreateCampaignInput} "Validation Error"
// @Failure 404 {object} response.Response{meta=response.Meta} "Campaign ID Not Found"
// @Failure 500 {object} response.Response{meta=response.Meta} "Internal Server Error"
func (r *Rest) UpdateCampaign(ctx *fiber.Ctx) error {
	var input request.UpdateCampaignInput
	var campaignID request.GetCampaignByIDInput

	err := ctx.ParamsParser(&campaignID)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	if errMap := validator.Validate(campaignID); errMap != nil {
		err := errors.NewErrorf(http.StatusBadRequest, errMap, "validation error")
		return response.RenderErrorResponse(ctx, "validation error", err)
	}

	err = ctx.BodyParser(&input)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	if errMap := validator.Validate(input); errMap != nil {
		err := errors.NewErrorf(http.StatusBadRequest, errMap, "validation error")
		return response.RenderErrorResponse(ctx, "validation error", err)
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
		return response.RenderErrorResponse(ctx, err.Error(), err)
	}

	return response.SuccessResponse(ctx, http.StatusOK, "OK", data)
}

func (r *Rest) UploadCampaignImage(ctx *fiber.Ctx) error {
	var input request.UploadCampaignImageInput
	err := ctx.BodyParser(&input)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	if errMap := validator.Validate(input); errMap != nil {
		err := errors.NewErrorf(http.StatusBadRequest, errMap, "validation error")
		return response.RenderErrorResponse(ctx, "validation error", err)
	}

	userResp, ok := ctx.Locals("current_user").(response.UserResponse)
	if !ok {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", nil)
	}

	input.User.ID = userResp.ID
	input.User.Email = userResp.Email

	imageFile, err := ctx.FormFile("campaign_image")

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err)
	}

	path := fmt.Sprintf("assets/images/campaigns/%s-%d-%s", input.CampaignID, time.Now().Day(), imageFile.Filename)

	err = ctx.SaveFile(imageFile, path)

	if err != nil {
		data := fiber.Map{
			"is_uploaded": false,
		}
		apiResponse := response.APIResponse("failed to upload campaign image", http.StatusBadRequest, "error", data, err.Error())
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
		res := response.APIResponse("failed to upload campaign image", http.StatusBadRequest, "error", data, err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(res)
	}

	return response.SuccessResponse(ctx, http.StatusOK, "OK", payload)
}
