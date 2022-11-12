package handler

import (
	"backend-crowdfunding/helper"
	"backend-crowdfunding/transaction"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type trxHandler struct {
	trxService transaction.Service
}

func NewTransactionHandler(transactionService transaction.Service) *trxHandler {
	return &trxHandler{transactionService}
}

func (h *trxHandler) GetAllTransactionsByCampaignID(ctx *gin.Context) {
	ID64int, err := strconv.ParseUint(ctx.Query("campaign_id"), 32, 64)
	campaignID := uint(ID64int)

	if err != nil {
		res := helper.APIresponse("Bad Request", 400, "error", nil, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	campaign, err := h.trxService.GetTransactionsByCampaignID(campaignID)
	if err != nil {
		res := helper.APIresponse("Bad Request", 400, "error", nil, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	response := helper.APIresponse("OK", http.StatusOK, "success", campaign, nil)
	ctx.JSON(http.StatusOK, response)
}
