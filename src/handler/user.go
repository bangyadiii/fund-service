package handler

import (
	"backend-crowdfunding/helper"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/request"
	"backend-crowdfunding/src/response"
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
		data := helper.APIResponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, errors)
		c.JSON(http.StatusBadRequest, data)
		return
	}
	var checkEmailFormatInput request.CheckEmailInput
	checkEmailFormatInput.Email = input.Email

	isAvailableEmail, err := r.service.User.IsEmailAvailable(c.Request.Context(), checkEmailFormatInput)

	if err != nil {
		data := helper.APIResponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}
	if !isAvailableEmail {
		respData := gin.H{"is_available": isAvailableEmail}
		data := helper.APIResponse("Email has been registered", http.StatusBadRequest, "error", respData, nil)
		c.JSON(http.StatusBadRequest, data)
		return
	}

	newUser, err := r.service.User.RegisterUser(c.Request.Context(), input)
	if err != nil {
		data := helper.APIResponse("Login failed.", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}
	newToken, err := r.service.Auth.GenerateToken(newUser.ID, newUser.Email)

	if err != nil {
		data := helper.APIResponse("Login failed.", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}
	formattedUser := response.FormatUserLogin(&newUser, newToken)

	if err != nil {
		data := helper.APIResponse("Register failed", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}

	data := helper.APIResponse("Register success", http.StatusOK, "success", formattedUser, nil)
	c.JSON(http.StatusOK, data)

}

func (r *rest) Login(c *gin.Context) {
	var input request.LoginUserInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatErrorValidation(err)
		helper.ErrorResponse(c, http.StatusBadRequest, "BAD REQUEST", errors)
		return
	}
	loginResponse, err := r.service.User.Login(c.Request.Context(), input)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "BAD REQUEST", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "OK", loginResponse)
}

func (r *rest) CheckIsEmailAvailable(c *gin.Context) {
	var input request.CheckEmailInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatErrorValidation(err)
		data := helper.APIResponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, errors)
		c.JSON(http.StatusUnprocessableEntity, data)
		return
	}
	IsEmailAvailable, err := r.service.User.IsEmailAvailable(c.Request.Context(), input)
	if err != nil {
		errors := helper.FormatErrorValidation(err)
		helper.ErrorResponse(c, http.StatusBadRequest, "BAD REQUEST", errors)
		return
	}
	respData := gin.H{"is_available": IsEmailAvailable}
	helper.SuccessResponse(c, http.StatusOK, "OK", respData)
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
		resp := helper.APIResponse("Failed to upload avatar 1", http.StatusBadRequest, "error", data, errors)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	filePath := fmt.Sprintf("assets/images/avatars/%s-%d%s", userID, time.Now().Unix(), path.Ext(file.Filename))
	err = c.SaveUploadedFile(file, filePath)

	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		resp := helper.APIResponse("Failed to upload avatar 2", http.StatusBadRequest, "error", data, err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	_, err = r.service.User.SaveAvatar(c.Request.Context(), userID, filePath)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		errors := helper.FormatErrorValidation(err)
		resp := helper.APIResponse("Failed to upload avatar 3", http.StatusBadRequest, "error", data, errors)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	data := gin.H{
		"is_uploaded": true,
	}
	resp := helper.APIResponse("Avatar successfuly uploaded.", http.StatusOK, "success", data, nil)
	c.JSON(http.StatusOK, resp)

}

func (r *rest) LoginWithGoogle(c *gin.Context) {
	var paramGoogle request.LoginWithGoogleInput
	err := c.ShouldBindJSON(&paramGoogle)
	if err != nil {
		data := helper.APIResponse("Login failed.", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}
	userRes, err := r.service.User.LoginWithGoogle(c.Request.Context(), paramGoogle)
	if err != nil {
		resp := helper.APIResponse("Failed to upload avatar 2", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	res := helper.APIResponse("Login success", http.StatusOK, "success", userRes, nil)
	c.JSON(http.StatusOK, res)
}
