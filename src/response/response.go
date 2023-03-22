package response

import (
	"math"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Meta       *Meta            `json:"meta"`
	Data       interface{}      `json:"data"`
	Pagination *PaginationParam `json:"pagination,omitempty"`
	Errors     interface{}      `json:"errors"`
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

func AddPagination(jsonResponse Response, pagination *PaginationParam) Response {
	jsonResponse.Pagination = pagination

	return jsonResponse
}

type PaginationParam struct {
	Limit        int64 `json:"limit" form:"limit" gorm:"-"`
	Offset       int64 `json:"offset" form:"-" gorm:"-"`
	Page         int64 `json:"page" form:"page" gorm:"-"`
	TotalPage    int64 `json:"total_page" gorm:"-"`
	CurrentPage  int64 `json:"current_page" gorm:"-"`
	TotalElement int64 `json:"total_element"`
}

func (pg *PaginationParam) ProcessPagination() bool {
	pg.CurrentPage = pg.Page
	pg.TotalPage = int64(math.Ceil(float64(pg.TotalElement) / float64(pg.Limit)))
	if pg.Page > pg.TotalPage {
		return false
	}

	pg.Offset = (pg.Page - 1) * pg.Limit
	return true
}

func FormatPaginationParam(params PaginationParam) *PaginationParam {
	paginationParam := params
	if params.Limit <= 0 {
		paginationParam.Limit = 10
	}
	if params.Page <= 0 {
		paginationParam.Page = 1
	}

	return &paginationParam
}

func SuccessResponseWithPagination(ctx *gin.Context, code int, message string, data interface{}, pg *PaginationParam) {
	json := APIResponse(message, code, "success", data, nil)
	json = AddPagination(json, pg)
	ctx.JSON(code, json)
}

func SuccessResponse(ctx *gin.Context, code int, message string, data interface{}) {
	json := APIResponse(message, code, "success", data, nil)
	ctx.JSON(code, json)
}

func ErrorResponse(ctx *gin.Context, code int, message string, errors interface{}) {
	json := APIResponse(message, code, "success", nil, errors)
	ctx.JSON(code, json)
}
