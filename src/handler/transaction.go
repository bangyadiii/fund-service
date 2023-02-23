package handler

import (
	"backend-crowdfunding/helper"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/request"
	"backend-crowdfunding/src/service"
	"net/http"

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

func (r *rest) GetAllTransactionsByCampaignID(ctx *gin.Context) {
	campaignID := ctx.Query("campaign_id")

	campaign, err := r.service.Trx.GetTransactionsByCampaignID(ctx.Request.Context(), campaignID)
	if err != nil {
		res := helper.APIresponse("Bad Request", 400, "error", nil, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	response := helper.APIresponse("OK", http.StatusOK, "success", campaign, nil)
	ctx.JSON(http.StatusOK, response)
}

func (r *rest) CreateTransaction(ctx *gin.Context) {
	var input request.CreateTransactionInput

	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		erros := helper.FormatErrorValidation(err)
		res := helper.APIresponse("Bad Request", http.StatusBadRequest, "error", nil, erros)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	input.UserID = ctx.MustGet("current_user").(model.User).ID

	trx, err := r.service.Trx.CreateTransaction(ctx.Request.Context(), input)
	if err != nil {
		res := helper.APIresponse("Bad Request", http.StatusBadRequest, "error", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.APIresponse("Created", http.StatusCreated, "success", trx, nil)
	ctx.JSON(http.StatusCreated, res)

}
