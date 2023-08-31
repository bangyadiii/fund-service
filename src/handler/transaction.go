package handler

import (
	"backend-crowdfunding/helper"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/request"
	"backend-crowdfunding/src/response"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type TransactionHandler interface {
	GetAllTransactionsByCampaignID(c *fiber.Ctx) error
	CreateTransaction(c *fiber.Ctx) error
}

// GetAllTransactionsByCampaignID A function that will be called when the user access the route `/transaction/campaign` with the method `GET`.
func (r *Rest) GetAllTransactionsByCampaignID(ctx *fiber.Ctx) error {
	campaignID := ctx.Query("campaign_id")

	campaignResp, err := r.service.Trx.GetTransactionsByCampaignID(ctx.Context(), campaignID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return response.SuccessResponse(ctx, http.StatusOK, "OK", campaignResp)
}

func (r *Rest) CreateTransaction(ctx *fiber.Ctx) error {
	var input request.CreateTransactionInput

	err := ctx.BodyParser(&input)

	if err != nil {
		errors := helper.FormatErrorValidation(err)
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", errors)
	}

	input.UserID = ctx.Locals("current_user").(model.User).ID

	trx, err := r.service.Trx.CreateTransaction(ctx.Context(), input)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return response.SuccessResponse(ctx, http.StatusCreated, "CREATED", trx)
}
