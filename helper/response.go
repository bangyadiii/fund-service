package helper

import "github.com/gin-gonic/gin"

type Response struct {
	Meta   *Meta       `json:"meta"`
	Data   interface{} `json:"data"`
	Errors interface{} `json:"errors"`
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

func SuccessResponse(ctx *gin.Context, code int, message string, data interface{}) {
	json := APIResponse(message, code, "success", data, nil)
	ctx.JSON(code, json)
}

func ErrorResponse(ctx *gin.Context, code int, message string, errors interface{}) {
	json := APIResponse(message, code, "success", nil, errors)
	ctx.JSON(code, json)
}
