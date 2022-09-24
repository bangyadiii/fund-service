package handler

import (
	"backend-crowdfunding/helper"
	"backend-crowdfunding/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)


type userHandler struct {
	userService user.Service
}
func NewUserHanlder(userService user.Service) *userHandler{
	return &userHandler{userService}	
}

func (h *userHandler) RegisterUser(c *gin.Context){
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		var errors []string

		for _, e := range err.(validator.ValidationErrors){
			errors = append(errors, e.Error())	
		}

		data :=  helper.APIresponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, errors)
		c.JSON(http.StatusBadRequest, data)
		return

	}
	newUser, err := h.userService.RegisterUser(input)
	formatedUser := user.FormatUser(newUser, "iniceritanya_token")

	if err != nil {
		data :=  helper.APIresponse("Register failed", http.StatusBadRequest, "error", nil,err)
		c.JSON(http.StatusBadRequest, data)
		return
	}

	data :=  helper.APIresponse("Register success", http.StatusOK, "success", formatedUser, nil)
	c.JSON(http.StatusOK, data)
	
}


func (h *userHandler) Login(c *gin.Context){
	var input user.LoginUserInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors){
			errors = append(errors, e.Error())
		}
		data := helper.APIresponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, errors)
		c.JSON(http.StatusUnprocessableEntity, data)
		return
	}
	validUser, err := h.userService.Login(input)
	if err != nil {
		data :=  helper.APIresponse("Login failed", http.StatusBadRequest, "error", nil, err)
		c.JSON(http.StatusBadRequest, data)
		return
	}

	formatedUser := user.FormatUser(validUser, "iniceritanya_token")
	data :=  helper.APIresponse("Login success", http.StatusOK, "success", formatedUser, nil)
	c.JSON(http.StatusOK, data)
		
}