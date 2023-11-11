package handler

import (
	"backend-crowdfunding/sdk/errors"
	"backend-crowdfunding/src/dto/request"
	"backend-crowdfunding/src/dto/response"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/validator"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type TransactionHandler interface {
	GetAllTransactionsByCampaignID(c *fiber.Ctx) error
	CreateTransaction(c *fiber.Ctx) error
}

// GetAllTransactionsByCampaignID A function that will be called when the user access the route `/transaction/campaign` with the method `GET`.
// @Summary Get Transaction.
// @Tags Transaction
// @Accept application/json
// @Produce json
// @Param Authorization header string true "Access token"
// @Param campaign_id query string true "Campaign ID"
// @Router /transactions/ [get]
// @Success 200 {object} response.Response{meta=response.Meta,data=[]response.TransactionResponse} "Success"
// @Failure 404 {object} response.Response{meta=response.Meta} "Transaction Not Found"
// @Failure 500 {object} response.Response{meta=response.Meta} "Internal Server Error
func (r *Rest) GetAllTransactionsByCampaignID(ctx *fiber.Ctx) error {
	campaignID := ctx.Query("campaign_id")

	campaignResp, err := r.service.Trx.GetTransactionsByCampaignID(ctx.Context(), campaignID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return response.SuccessResponse(ctx, http.StatusOK, "OK", campaignResp)
}

// CreateTransaction handling transaction creation.
// @Summary Create Transaction.
// @Tags Transaction
// @Accept application/json
// @Produce json
// @Param Authorization header string true "Access token"
// @Param request body request.CreateTransactionInput true "Transaction payload"
// @Router /transactions/ [post]
// @Success 201 {object} response.Response{meta=response.Meta,data=response.TransactionResponse} "Success"
// @Failure 400 {object} response.Response{meta=response.Meta,errors=request.CreateTransactionInput} "Validation Error"
// @Failure 404 {object} response.Response{meta=response.Meta} "Transaction Not Found"
// @Failure 500 {object} response.Response{meta=response.Meta} "Internal Server Error
func (r *Rest) CreateTransaction(ctx *fiber.Ctx) error {
	var input request.CreateTransactionInput

	err := ctx.BodyParser(&input)

	if err != nil {
		return response.RenderErrorResponse(ctx, "BAD REQUEST", err)
	}

	if errMap := validator.Validate(input); errMap != nil {
		err := errors.NewErrorf(http.StatusBadRequest, errMap, "validation error")
		return response.RenderErrorResponse(ctx, "validation error", err)
	}

	input.UserID = ctx.Locals("current_user").(model.User).ID

	trx, err := r.service.Trx.CreateTransaction(ctx.Context(), input)
	if err != nil {
		return response.RenderErrorResponse(ctx, err.Error(), err)
	}

	return response.SuccessResponse(ctx, http.StatusCreated, "CREATED", trx)
}
