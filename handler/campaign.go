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
		response := helper.APIresponse("Error occur while getting campaing", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	res := helper.APIresponse("OK", http.StatusOK, "success", user, nil)

	c.JSON(http.StatusOK, res)
}
