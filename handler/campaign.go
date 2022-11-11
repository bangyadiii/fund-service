package handler

import (
	"backend-crowdfunding/campaign"
	"backend-crowdfunding/helper"
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
	user_id, _ := strconv.Atoi(c.Query("user_id"))

	user, err := h.campaignService.GetCampaigns(user_id)

	if err != nil {
		response := helper.APIresponse("Error occur while getting campaign", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	res := helper.APIresponse("OK", http.StatusOK, "success", user, nil)

	c.JSON(http.StatusOK, res)
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

	data, err := h.campaignService.CreateCampaign(input)
	if err != nil {
		res := helper.APIresponse("Something went wrong", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}
	formattedCampaign := campaign.FormatCampaign(data)
	res := helper.APIresponse("CREATED", http.StatusCreated, "success", formattedCampaign, nil)
	c.JSON(http.StatusCreated, res)
}
