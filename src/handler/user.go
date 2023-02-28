package handler

import (
	"backend-crowdfunding/helper"
	"backend-crowdfunding/src/formatter"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/request"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	RegisterUser(c *gin.Context)
	Login(c *gin.Context)
	CheckIsEmailAvailable(c *gin.Context)
	UploadAvatar(c *gin.Context)
}

func (r *rest) RegisterUser(c *gin.Context) {
	var input request.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatErrorValidation(err)
		data := helper.APIresponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, errors)
		c.JSON(http.StatusBadRequest, data)
		return
	}
	var checkEmailFormatInput request.CheckEmailInput
	checkEmailFormatInput.Email = input.Email

	isAvailableEmail, err := r.service.User.IsEmailAvailable(c.Request.Context(), checkEmailFormatInput)

	if err != nil {
		data := helper.APIresponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}
	if !isAvailableEmail {
		respData := gin.H{"is_available": isAvailableEmail}
		data := helper.APIresponse("Email has been registered", http.StatusBadRequest, "error", respData, nil)
		c.JSON(http.StatusBadRequest, data)
		return
	}

	newUser, err := r.service.User.RegisterUser(c.Request.Context(), input)
	if err != nil {
		data := helper.APIresponse("Login failed.", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}
	newToken, err := r.service.Auth.GenerateToken(newUser.ID, newUser.Email)

	if err != nil {
		data := helper.APIresponse("Login failed.", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}
	formatedUser := formatter.FormatUserLogin(newUser, newToken)

	if err != nil {
		data := helper.APIresponse("Register failed", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}

	data := helper.APIresponse("Register success", http.StatusOK, "success", formatedUser, nil)
	c.JSON(http.StatusOK, data)

}

func (r *rest) Login(c *gin.Context) {
	var input request.LoginUserInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatErrorValidation(err)
		data := helper.APIresponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, errors)
		c.JSON(http.StatusUnprocessableEntity, data)
		return
	}
	validUser, err := r.service.User.Login(c.Request.Context(), input)
	if err != nil {

		data := helper.APIresponse("Login failed", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}
	newToken, err := r.service.Auth.GenerateToken(validUser.ID, validUser.Email)

	if err != nil {
		data := helper.APIresponse("Login failed.", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}

	formatedUser := formatter.FormatUserLogin(validUser, newToken)
	data := helper.APIresponse("Login success", http.StatusOK, "success", formatedUser, nil)
	c.JSON(http.StatusOK, data)

}

func (r *rest) CheckIsEmailAvailable(c *gin.Context) {
	var input request.CheckEmailInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatErrorValidation(err)
		data := helper.APIresponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, errors)
		c.JSON(http.StatusUnprocessableEntity, data)
		return
	}
	IsEmailAvailable, err := r.service.User.IsEmailAvailable(c.Request.Context(), input)
	if err != nil {
		errors := helper.FormatErrorValidation(err)
		data := helper.APIresponse("Error", http.StatusBadRequest, "error", nil, errors)
		c.JSON(http.StatusBadRequest, data)
		return
	}
	if !IsEmailAvailable {
		respData := gin.H{"is_available": IsEmailAvailable}
		data := helper.APIresponse("Email has been registered", http.StatusOK, "error", respData, nil)
		c.JSON(http.StatusOK, data)
		return
	}
	respData := gin.H{"is_available": IsEmailAvailable}
	data := helper.APIresponse("Email is available", http.StatusOK, "success", respData, nil)
	c.JSON(http.StatusOK, data)

}

func (r *rest) UploadAvatar(c *gin.Context) {

	currentUser := c.MustGet("current_user").(model.User)
	userID := currentUser.ID
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		errors := helper.FormatErrorValidation(err)
		response := helper.APIresponse("Failed to upload avatar 1", http.StatusBadRequest, "error", data, errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	filePath := fmt.Sprintf("assets/images/avatars/%s-%d%s", userID, time.Now().Unix(), path.Ext(file.Filename))
	err = c.SaveUploadedFile(file, filePath)

	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIresponse("Failed to upload avatar 2", http.StatusBadRequest, "error", data, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = r.service.User.SaveAvatar(c.Request.Context(), userID, filePath)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		errors := helper.FormatErrorValidation(err)
		response := helper.APIresponse("Failed to upload avatar 3", http.StatusBadRequest, "error", data, errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_uploaded": true,
	}
	response := helper.APIresponse("Avatar successfuly uploaded.", http.StatusOK, "success", data, nil)
	c.JSON(http.StatusOK, response)

}

func (r *rest) LoginWithGoogle(c *gin.Context) {
	var paramGoogle request.LoginWithGoogleInput
	err := c.ShouldBindJSON(&paramGoogle)
	if err != nil {
		data := helper.APIresponse("Login failed.", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}
	userRes, err := r.service.User.LoginWithGoogle(c.Request.Context(), paramGoogle)
	if err != nil {
		response := helper.APIresponse("Failed to upload avatar 2", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	res := helper.APIresponse("Login success", http.StatusOK, "success", userRes, nil)
	c.JSON(http.StatusOK, res)
}
