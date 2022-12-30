package handler

import (
	"backend-crowdfunding/helper"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/request"
	"backend-crowdfunding/src/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionHandler interface {
	GetAllTransactionsByCampaignID(c *gin.Context)
	CreateTransaction(c *gin.Context)
}

type trxHandler struct {
	trxService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *trxHandler {
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

func (h *trxHandler) CreateTransaction(ctx *gin.Context) {
	var input request.CreateTransactionInput

	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		erros := helper.FormatErrorValidation(err)
		res := helper.APIresponse("Bad Request", http.StatusBadRequest, "error", nil, erros)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	input.UserID = ctx.MustGet("current_user").(model.User).ID

	trx, err := h.trxService.CreateTransaction(input)
	if err != nil {
		res := helper.APIresponse("Bad Request", http.StatusBadRequest, "error", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.APIresponse("Created", http.StatusCreated, "success", trx, nil)
	ctx.JSON(http.StatusCreated, res)

}
