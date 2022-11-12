package handler

import (
	"backend-crowdfunding/campaign"
	"backend-crowdfunding/helper"
	"backend-crowdfunding/user"
	"net/http"
	"strconv"

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
