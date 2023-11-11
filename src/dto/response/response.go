package response

import (
	ierrors "backend-crowdfunding/sdk/errors"
	"backend-crowdfunding/src/dto/request"
	"errors"
	"math"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Meta   *Meta       `json:"meta"`
	Data   interface{} `json:"data"`
	Errors interface{} `json:"errors"`
}

type WithPagination struct {
	Meta       *Meta               `json:"meta"`
	Data       interface{}         `json:"data"`
	Pagination *PaginationResponse `json:"pagination,omitempty" extensions:"x-nullable,x-omitempty"`
	Errors     interface{}         `json:"errors" extensions:"x-nullable"`
}

type Meta struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func APIResponse(message string, code int, status string, data interface{}, errors interface{}) Response {
	metaData := Meta{
		Message: message,
		Status:  status,
		Code:    code,
	}

	jsonResponse := Response{
		Meta:   &metaData,
		Data:   data,
		Errors: errors,
	}

	return jsonResponse
}

func RenderErrorResponse(ctx *fiber.Ctx, message string, err error) error {
	var errorMap map[string]string
	var status = http.StatusInternalServerError
	var ierr *ierrors.Error
	if !errors.As(err, &ierr) {
		message = "Internal Server Error"
	} else {
		message = ierr.Error()
		switch ierr.GetCode() {
		case ierrors.NotFoundType:
			status = http.StatusNotFound
		case ierrors.ValidationErrorType:
			status = http.StatusBadRequest
			errorMap = ierr.Errors()
		}
	}

	return ErrorResponse(ctx, status, message, errorMap)
}

func AddPagination(jsonResponse Response, pagination *PaginationResponse) WithPagination {
	jsonPagination := WithPagination{
		Meta:       jsonResponse.Meta,
		Data:       jsonResponse.Data,
		Errors:     jsonResponse.Errors,
		Pagination: pagination,
	}

	return jsonPagination
}

type PaginationResponse struct {
	Limit        int64 `json:"limit" gorm:"-"`
	Offset       int64 `json:"offset" gorm:"-"`
	Page         int64 `json:"page" gorm:"-"`
	TotalPage    int64 `json:"total_page" gorm:"-" form:"-"`
	CurrentPage  int64 `json:"current_page" gorm:"-"`
	TotalElement int64 `json:"total_element"`
}

func ConvertPaginationParamToPaginationResponse(req request.PaginationParam) *PaginationResponse {
	return &PaginationResponse{
		Limit:        req.Limit,
		Offset:       0,
		Page:         req.Page,
		TotalPage:    0,
		CurrentPage:  0,
		TotalElement: 0,
	}
}

func (pg *PaginationResponse) ProcessPagination() bool {
	pg.CurrentPage = pg.Page
	pg.TotalPage = int64(math.Ceil(float64(pg.TotalElement) / float64(pg.Limit)))
	if pg.Page > pg.TotalPage {
		return false
	}

	pg.Offset = (pg.Page - 1) * pg.Limit
	return true
}

func FormatPaginationParam(params PaginationResponse) *PaginationResponse {
	paginationParam := params
	if params.Limit <= 0 {
		paginationParam.Limit = 10
	}
	if params.Page <= 0 {
		paginationParam.Page = 1
	}

	return &paginationParam
}

func SuccessResponseWithPagination(ctx *fiber.Ctx, code int, message string, data interface{}, pg *PaginationResponse) error {
	json := APIResponse(message, code, "success", data, nil)
	paginationResp := AddPagination(json, pg)
	return ctx.Status(code).JSON(paginationResp)
}

func SuccessResponse(ctx *fiber.Ctx, code int, message string, data interface{}) error {

	json := APIResponse(message, code, statusMap[code], data, nil)
	return ctx.Status(code).JSON(json)
}

func ErrorResponse(ctx *fiber.Ctx, code int, message string, errors interface{}) error {
	json := APIResponse(message, code, "error", nil, errors)
	return ctx.Status(code).JSON(json)
}

var statusMap = map[int]string{
	200: "OK",
	201: "Created",
	202: "Accepted",
	204: "No Content",
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	409: "Conflict",
	422: "Unprocessable Entity",
	500: "Internal Server Error",
	501: "Not Implemented",
	503: "Service Unavailable",
}
