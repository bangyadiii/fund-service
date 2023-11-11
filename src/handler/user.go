package handler

import (
	"backend-crowdfunding/sdk/errors"
	"backend-crowdfunding/src/dto/request"
	"backend-crowdfunding/src/dto/response"
	"backend-crowdfunding/src/validator"
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

// RegisterUser Register User
// @Summary Register new user
// @Tags Authentication
// @Accept application/json
// @Produce json
// @Success 200 {object} response.Response{meta=response.Meta,data=response.UserLoginResponse} "Success"
// @Failure 400 {object} response.Response{meta=response.Meta,errors=request.RegisterUserInput} "Validation Error"
// @Failure 409 {object} response.Response{meta=response.Meta,errors=map[string]string} "Conflict email"
// @Router /auth/register [post]
func (r *Rest) RegisterUser(ctx *fiber.Ctx) error {
	var input = new(request.RegisterUserInput)

	err := ctx.BodyParser(input)
	if err != nil {
		return response.RenderErrorResponse(ctx, "Bad Request", err)
	}

	if errMap := validator.Validate(input); errMap != nil {
		err := errors.NewErrorf(http.StatusBadRequest, errMap, "Validation Error")
		return response.RenderErrorResponse(ctx, err.Error(), err)
	}

	newUser, newToken, err := r.service.User.RegisterUser(ctx.Context(), *input)

	if err != nil {
		return response.RenderErrorResponse(ctx, "Register failed", err)
	}

	formattedUser := response.FormatUserLogin(&newUser, newToken)
	return response.SuccessResponse(ctx, http.StatusOK, "success", formattedUser)
}

// Login log in user
// @Summary Login User
// @Tags Authentication
// @Accept application/json
// @Produce json
// @Success 200 {object} response.Response{meta=response.Meta,data=response.UserLoginResponse} "Success"
// @Failure 400 {object} response.Response{meta=response.Meta,errors=request.LoginUserInput} "Validation Error"
// @Router /auth/login [post]
func (r *Rest) Login(ctx *fiber.Ctx) error {
	var input = new(request.LoginUserInput)
	err := ctx.BodyParser(&input)

	if err != nil {
		return response.RenderErrorResponse(ctx, "Bad Request", err)
	}

	if errMap := validator.Validate(input); errMap != nil {
		err := errors.NewErrorf(http.StatusBadRequest, errMap, "validation error")
		return response.RenderErrorResponse(ctx, "validation error", err)
	}

	loginResponse, err := r.service.User.Login(ctx.Context(), *input)

	if err != nil {
		return response.RenderErrorResponse(ctx, "Login Failed", err)
	}

	return response.SuccessResponse(ctx, http.StatusOK, "OK", loginResponse)
}

// CheckIsEmailAvailable log in user
// @Summary Check is email available
// @Tags Authentication
// @Accept application/json
// @Produce json
// @Router /auth/email-is-available [post]
// @Success 200 {object} response.Response{data=map[string]string} "Success"
// @Failure 400 {object} response.Response{errors=request.CheckEmailInput} "Validation Error"
func (r *Rest) CheckIsEmailAvailable(ctx *fiber.Ctx) error {
	var input request.CheckEmailInput
	err := ctx.BodyParser(&input)

	if err != nil {
		return response.RenderErrorResponse(ctx, "Bad Request", err)
	}

	if errMap := validator.Validate(input); errMap != nil {
		err := errors.NewErrorf(http.StatusBadRequest, errMap, "validation error")
		return response.RenderErrorResponse(ctx, "validation error", err)
	}

	IsEmailAvailable, err := r.service.User.IsEmailAvailable(ctx.Context(), input)
	if err != nil {
		return response.RenderErrorResponse(ctx, err.Error(), err)
	}

	return response.SuccessResponse(ctx, http.StatusOK, "OK", fiber.Map{"is_available": IsEmailAvailable})
}

func (r *Rest) UploadAvatar(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals("current_user").(response.UserResponse)
	userID := currentUser.ID
	file, err := ctx.FormFile("avatar")
	if err != nil {
		return response.RenderErrorResponse(ctx, "Bad Request", err)
	}

	// TODO: create storage layer that implement upload file
	filePath := fmt.Sprintf("assets/images/avatars/%s-%d%s", userID, time.Now().Unix(), path.Ext(file.Filename))
	err = ctx.SaveFile(file, filePath)

	if err != nil {
		err := errors.NewErrorf(400, nil, "failed to upload avatar", err.Error())
		return response.RenderErrorResponse(ctx, err.Error(), err)
	}

	_, err = r.service.User.SaveAvatar(ctx.Context(), userID, filePath)
	if err != nil {
		return response.RenderErrorResponse(ctx, err.Error(), err)
	}

	// TODO: instead of returning `is_uploaded: true`, better to returns its public URL
	data := fiber.Map{
		"is_uploaded": true,
	}

	return response.SuccessResponse(ctx, http.StatusOK, "Avatar successfully uploaded.", data)
}

func (r *Rest) LoginWithGoogle(ctx *fiber.Ctx) error {
	var paramGoogle = new(request.LoginWithGoogleInput)
	err := ctx.BodyParser(paramGoogle)

	if err != nil {
		return response.RenderErrorResponse(ctx, "Login failed", err)
	}

	if errMap := validator.Validate(paramGoogle); errMap != nil {
		err := errors.NewErrorf(http.StatusBadRequest, errMap, "validation error")
		return response.RenderErrorResponse(ctx, "validation error", err)
	}

	userRes, err := r.service.User.LoginWithGoogle(ctx.Context(), *paramGoogle)
	if err != nil {
		err := errors.NewErrorf(400, nil, "failed to upload avatar", err.Error())
		return response.RenderErrorResponse(ctx, "Bad Request", err)
	}

	return response.SuccessResponse(ctx, http.StatusOK, "Login success", userRes)
}
