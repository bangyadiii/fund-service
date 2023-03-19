package handler

import (
	"backend-crowdfunding/helper"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler interface {
	GetAllTransactionsByCampaignID(c *gin.Context)
	CreateTransaction(c *gin.Context)
}

// GetAllTransactionsByCampaignID A function that will be called when the user access the route `/transaction/campaign` with the method `GET`.
func (r *rest) GetAllTransactionsByCampaignID(ctx *gin.Context) {
	campaignID := ctx.Query("campaign_id")

	campaignResp, err := r.service.Trx.GetTransactionsByCampaignID(ctx.Request.Context(), campaignID)
	if err != nil {
		helper.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
		return
	}

	helper.SuccessResponse(ctx, http.StatusOK, "OK", campaignResp)
}

func (r *rest) CreateTransaction(ctx *gin.Context) {
	var input request.CreateTransactionInput

	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatErrorValidation(err)
		helper.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", errors)
		return
	}

	input.UserID = ctx.MustGet("current_user").(model.User).ID

	trx, err := r.service.Trx.CreateTransaction(ctx.Request.Context(), input)
	if err != nil {
		helper.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
		return
	}

	helper.SuccessResponse(ctx, http.StatusCreated, "CREATED", trx)
}
