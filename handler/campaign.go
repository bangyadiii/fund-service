package handler

import (
	"backend-crowdfunding/campaign"
	"backend-crowdfunding/helper"
	"backend-crowdfunding/user"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	user_id, _ := strconv.ParseUint(c.Query("user_id"), 32, 64)
	userID := uint(user_id)
	data, err := h.campaignService.GetCampaigns(userID)

	if err != nil {
		response := helper.APIresponse("Error occur while getting campaign", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	payload := campaign.FormatCampaignCollections(data)

	res := helper.APIresponse("OK", http.StatusOK, "success", payload, nil)

	c.JSON(http.StatusOK, res)
}

func (h *campaignHandler) GetCampaignByID(ctx *gin.Context) {
	var input campaign.GetCampaignByIDInput
	err := ctx.ShouldBindUri(&input)

	if err != nil {
		res := helper.APIresponse("Bad request", http.StatusBadRequest, "error", nil, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := h.campaignService.GetCampaignByID(input)

	if err != nil {
		res := helper.APIresponse("Something went wrong", http.StatusInternalServerError, "error", nil, err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	formatted := campaign.FormatCampaignDetail(data)
	payload := helper.APIresponse("OK", http.StatusOK, "success", formatted, nil)
	ctx.JSON(http.StatusOK, payload)
}

func (h *campaignHandler) CreateNewCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatErrorValidation(err)
		res := helper.APIresponse("Something went wrong", http.StatusBadRequest, "error", nil, errors)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	curUser := c.MustGet("current_user").(user.User)
	input.User = curUser

	data, err := h.campaignService.CreateCampaign(input)
	if err != nil {
		res := helper.APIresponse("Something went wrong", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}
	formattedCampaign := campaign.FormatCampaignDetail(data)
	res := helper.APIresponse("CREATED", http.StatusCreated, "success", formattedCampaign, nil)
	c.JSON(http.StatusCreated, res)
}

func (h *campaignHandler) UpdateCampaign(ctx *gin.Context) {
	var input campaign.UpdateCampaignInput
	var campaignID campaign.GetCampaignByIDInput

	err := ctx.ShouldBindUri(&campaignID)

	if err != nil {
		errors := helper.FormatErrorValidation(err)
		res := helper.APIresponse("Bad Request", http.StatusBadRequest, "error", nil, errors)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err = ctx.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatErrorValidation(err)
		res := helper.APIresponse("Bad Request", http.StatusBadRequest, "error", nil, errors)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	curUser := ctx.MustGet("current_user").(user.User)
	input.User = curUser

	data, err := h.campaignService.UpdateCampaign(campaignID, input)

	if err != nil {
		res := helper.APIresponse("Something went wrong", http.StatusInternalServerError, "error", nil, err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	formattedCampaign := campaign.FormatCampaignDetail(data)
	res := helper.APIresponse("OK", http.StatusOK, "success", formattedCampaign, nil)

	ctx.JSON(http.StatusOK, res)
}

func (h *campaignHandler) UploadCampaignImage(ctx *gin.Context) {
	var input campaign.UploadCampaignImageInput
	err := ctx.ShouldBind(&input)

	if err != nil {
		erros := helper.FormatErrorValidation(err)
		res := helper.APIresponse("Bad Reqeust", http.StatusBadRequest, "error", nil, erros)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	imageFile, err := ctx.FormFile("campaign_image")

	if err != nil {
		res := helper.APIresponse("Bad Reqeust", http.StatusBadRequest, "error", nil, err)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	path := fmt.Sprintf("assets/images/campaigns/%d-%d-%s", input.CampaignID, time.Now().Day(), imageFile.Filename)

	err = ctx.SaveUploadedFile(imageFile, path)

	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIresponse("Failed to upload campaign image", http.StatusBadRequest, "error", data, err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	input.ImageName = path
	payload, err := h.campaignService.UploadCampaignImage(input)

	if err != nil {
		os.Remove(path)
		data := gin.H{
			"is_uploaded": false,
		}
		res := helper.APIresponse("Failed to upload campaign image", http.StatusBadRequest, "error", data, err.Error())
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.APIresponse("Avatar successfuly uploaded.", http.StatusOK, "success", payload, nil)
	ctx.JSON(http.StatusOK, response)
}
