package middleware

import (
	"backend-crowdfunding/src/response"
	"backend-crowdfunding/src/service"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func VerifyToken(userService service.UserService, authService service.AuthService) fiber.Handler {
	/*
		This code is used to authenticate a user and authorize them to access a certain endpoint.
		It first checks if the Authorization header contains "Bearer". If not, it returns an error response.
		It then uses authService to validate the tokenString, and if there is an error, it returns an error response.
		It then checks if the claims are valid and if not, it returns an error response.
		Finally, it uses userService to find the user by ID and check if the email matches with what is stored in payload.
		If all checks pass, it sets a current_user variable with the user object.
	*/

	return func(ctx *fiber.Ctx) error {
		bearerToken := ctx.Get("Authorization")

		if !strings.Contains(bearerToken, "Bearer") {
			resp := response.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil, "Cannot hit this endpoint with no authentication. 1")
			return ctx.Status(http.StatusUnauthorized).JSON(resp)
		}
		splitBearer := strings.Split(bearerToken, " ")
		var tokenString string = ""
		if len(splitBearer) == 2 {
			tokenString = splitBearer[1]
		}
		accessToken, err := authService.ValidateToken(tokenString)
		if err != nil {
			resp := response.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil, err.Error())
			return ctx.Status(http.StatusUnauthorized).JSON(resp)
		}
		payload, ok := authService.ConvertTokenToCustomClaims(accessToken)

		if !ok || !accessToken.Valid {
			resp := response.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil, fiber.Map{"message": "Do not have permissions to access this resources"})
			return ctx.Status(http.StatusUnauthorized).JSON(resp)
		}

		userID := payload.ID
		user, err := userService.FindByID(ctx.Context(), userID)

		if err != nil {
			resp := response.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil, "Cannot hit this endpoint with no authentication.")
			return ctx.Status(http.StatusUnauthorized).JSON(resp)
		}

		if user.Email != payload.Email {
			resp := response.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil, "Cannot hit this endpoint with no authentication.")
			return ctx.Status(http.StatusUnauthorized).JSON(resp)
		}
		ctx.Locals("current_user", user)
		return nil
	}
}
