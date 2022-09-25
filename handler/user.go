package handler

import (
	"backend-crowdfunding/auth"
	"backend-crowdfunding/helper"
	"backend-crowdfunding/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)


type userHandler struct {
	userService user.Service
	authService auth.Service
}
func NewUserHanlder(userService user.Service, authService auth.Service) *userHandler{
	return &userHandler{userService, authService}	
}

func (h *userHandler) RegisterUser(c *gin.Context){
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatErrorValidation(err)
		data :=  helper.APIresponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, errors)
		c.JSON(http.StatusBadRequest, data)
		return

	}
	newUser, err := h.userService.RegisterUser(input)
	newToken, err := h.authService.GenerateToken(newUser.ID, newUser.Email)

	if err != nil {
		data :=  helper.APIresponse("Login failed.", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}
	formatedUser := user.FormatUser(newUser, newToken)

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
		errors := helper.FormatErrorValidation(err)	
		data := helper.APIresponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, errors)
		c.JSON(http.StatusUnprocessableEntity, data)
		return
	}
	validUser, err := h.userService.Login(input)
	if err != nil {

		data :=  helper.APIresponse("Login failed", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}
	newToken, err := h.authService.GenerateToken(validUser.ID, validUser.Email)

	if err != nil {
		data :=  helper.APIresponse("Login failed.", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, data)
		return
	}

	formatedUser := user.FormatUser(validUser,newToken)
	data :=  helper.APIresponse("Login success", http.StatusOK, "success", formatedUser, nil)
	c.JSON(http.StatusOK, data)
		
}


func (h *userHandler) CheckIsEmailAvailable(c *gin.Context){
	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatErrorValidation(err)	
		data := helper.APIresponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, errors)
		c.JSON(http.StatusUnprocessableEntity, data)
		return
	}
	IsEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errors := helper.FormatErrorValidation(err)
		data :=  helper.APIresponse("Error", http.StatusBadRequest, "error", nil, errors)
		c.JSON(http.StatusBadRequest, data)
		return
	}
	if !IsEmailAvailable {
		respData := gin.H{"is_available" : IsEmailAvailable}
		data :=  helper.APIresponse("Email has been registered", http.StatusOK, "error", respData, nil)
		c.JSON(http.StatusOK, data)
		return
	}
	respData := gin.H{"is_available" : IsEmailAvailable}
	data :=  helper.APIresponse("Email is available", http.StatusOK, "success", respData, nil)
	c.JSON(http.StatusOK, data)
	
}

func (h *userHandler) UploadAvatar(c *gin.Context){
	// TODO input dari user
	// TODO simpan gambarnya di folder /images/
	// TODO di service panggil repo untuk simpan nama gambar
	// TODO JWT (sementara hardcode, seakan2 user ID 1 yang akan upload avatar)
	// TODO repo ambil data user ID 1
	// TODO repo update data user dengan data nama avatar/ lokasi file avatar

	userID := 1
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		errors := helper.FormatErrorValidation(err)
		response :=  helper.APIresponse("Failed to upload avatar 1", http.StatusBadRequest, "error", data, errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	path := fmt.Sprintf("assets/images/%d-%d-%s", userID, file.Size, file.Filename) 
	err = c.SaveUploadedFile(file,path)
	
	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		response :=  helper.APIresponse("Failed to upload avatar 2", http.StatusBadRequest, "error", data, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID,path)
	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		errors := helper.FormatErrorValidation(err)
		response :=  helper.APIresponse("Failed to upload avatar 3", http.StatusBadRequest, "error", data, errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_uploaded" : true,
	}
	response :=  helper.APIresponse("Avatar successfuly uploaded.", http.StatusOK, "success", data, nil)
	c.JSON(http.StatusOK, response)

}



