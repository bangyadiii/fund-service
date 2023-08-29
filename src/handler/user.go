package handler

import (
	"backend-crowdfunding/helper"
	"backend-crowdfunding/src/request"
	"backend-crowdfunding/src/response"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	RegisterUser(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	CheckIsEmailAvailable(ctx *fiber.Ctx) error
	UploadAvatar(ctx *fiber.Ctx) error
}

func (r *rest) RegisterUser(ctx *fiber.Ctx) error {
	var input request.RegisterUserInput

	err := ctx.BodyParser(&input)
	if err != nil {
		errors := helper.FormatErrorValidation(err)
		data := response.APIResponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, errors)
		return ctx.Status(http.StatusBadRequest).JSON(data)
	}
	var checkEmailFormatInput request.CheckEmailInput
	checkEmailFormatInput.Email = input.Email

	isAvailableEmail, err := r.service.User.IsEmailAvailable(ctx.Context(), checkEmailFormatInput)

	if err != nil {
		data := response.APIResponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(data)
	}
	if !isAvailableEmail {
		respData := fiber.Map{"is_available": isAvailableEmail}
		data := response.APIResponse("Email has been registered", http.StatusBadRequest, "error", respData, nil)
		return ctx.Status(http.StatusBadRequest).JSON(data)
	}

	newUser, err := r.service.User.RegisterUser(ctx.Context(), input)
	if err != nil {
		data := response.APIResponse("Login failed.", http.StatusBadRequest, "error", nil, err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(data)
	}
	newToken, err := r.service.Auth.GenerateToken(newUser.ID, newUser.Email)

	if err != nil {
		data := response.APIResponse("Login failed.", http.StatusBadRequest, "error", nil, err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(data)
	}
	formattedUser := response.FormatUserLogin(&newUser, newToken)

	if err != nil {
		data := response.APIResponse("Register failed", http.StatusBadRequest, "error", nil, err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(data)
	}

	data := response.APIResponse("Register success", http.StatusOK, "success", formattedUser, nil)
	return ctx.Status(http.StatusOK).JSON(data)

}

func (r *rest) Login(ctx *fiber.Ctx) error {
	var input request.LoginUserInput
	err := ctx.BodyParser(&input)

	if err != nil {
		errors := helper.FormatErrorValidation(err)
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", errors)
	}
	loginResponse, err := r.service.User.Login(ctx.Context(), input)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return response.SuccessResponse(ctx, http.StatusOK, "OK", loginResponse)
}

func (r *rest) CheckIsEmailAvailable(c *fiber.Ctx) error {
	var input request.CheckEmailInput
	err := c.BodyParser(&input)

	if err != nil {
		errors := helper.FormatErrorValidation(err)
		data := response.APIResponse("Bad Request", http.StatusUnprocessableEntity, "error", nil, errors)
		return c.Status(http.StatusUnprocessableEntity).JSON(data)
	}
	IsEmailAvailable, err := r.service.User.IsEmailAvailable(c.Context(), input)
	if err != nil {
		return err

		// errors := helper.FormatErrorValidation(err)
		// response.ErrorResponse(c, http.StatusBadRequest, "BAD REQUEST", errors)
	}
	respData := fiber.Map{"is_available": IsEmailAvailable}
	return response.SuccessResponse(c, http.StatusOK, "OK", respData)
}

func (r *rest) UploadAvatar(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals("current_user").(response.UserResponse)
	userID := currentUser.ID
	file, err := ctx.FormFile("avatar")
	if err != nil {
		data := fiber.Map{
			"is_uploaded": false,
		}
		errors := helper.FormatErrorValidation(err)
		resp := response.APIResponse("Failed to upload avatar 1", http.StatusBadRequest, "error", data, errors)
		return ctx.Status(http.StatusBadRequest).JSON(resp)
	}

	filePath := fmt.Sprintf("assets/images/avatars/%s-%d%s", userID, time.Now().Unix(), path.Ext(file.Filename))
	err = ctx.SaveFile(file, filePath)

	if err != nil {
		data := fiber.Map{
			"is_uploaded": false,
		}
		resp := response.APIResponse("Failed to upload avatar 2", http.StatusBadRequest, "error", data, err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(resp)
	}

	_, err = r.service.User.SaveAvatar(ctx.Context(), userID, filePath)
	if err != nil {
		data := fiber.Map{
			"is_uploaded": false,
		}
		errors := helper.FormatErrorValidation(err)
		resp := response.APIResponse("Failed to upload avatar 3", http.StatusBadRequest, "error", data, errors)
		return ctx.Status(http.StatusBadRequest).JSON(resp)
	}

	data := fiber.Map{
		"is_uploaded": true,
	}
	resp := response.APIResponse("Avatar successfuly uploaded.", http.StatusOK, "success", data, nil)
	return ctx.Status(http.StatusOK).JSON(resp)
}

func (r *rest) LoginWithGoogle(ctx *fiber.Ctx) error {
	var paramGoogle request.LoginWithGoogleInput
	err := ctx.BodyParser(&paramGoogle)
	if err != nil {
		data := response.APIResponse("Login failed.", http.StatusBadRequest, "error", nil, err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(data)
	}
	userRes, err := r.service.User.LoginWithGoogle(ctx.Context(), paramGoogle)
	if err != nil {
		resp := response.APIResponse("Failed to upload avatar 2", http.StatusBadRequest, "error", nil, err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(resp)
	}
	res := response.APIResponse("Login success", http.StatusOK, "success", userRes, nil)
	return ctx.Status(http.StatusOK).JSON(res)
}
