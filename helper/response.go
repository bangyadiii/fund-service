package helper

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

func APIresponse(message string, code int, status string, data interface{}, errors interface{}) Response {
	meta_data := Meta{
		Message: message,
		Status:  status,
		Code:    code,
	}

	jsonResponse := Response{
		Meta:   &meta_data,
		Data:   data,
		Errors: errors,
	}

	return jsonResponse
}
