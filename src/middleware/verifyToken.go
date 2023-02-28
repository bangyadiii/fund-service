package middleware

import (
	"backend-crowdfunding/helper"
	"backend-crowdfunding/src/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func VerifyToken(userService service.UserService, authService service.AuthService) gin.HandlerFunc {
	/*
		This code is used to authenticate a user and authorize them to access a certain endpoint.
		It first checks if the Authorization header contains "Bearer". If not, it returns an error response.
		It then splits the header into two parts and stores the token string in tokenString.
		It then uses authService to validate the tokenString, and if there is an error, it returns an error response.
		It then checks if the claims are valid and if not, it returns an error response.
		Finally, it uses userService to find the user by ID and check if the email matches with what is stored in payload.
		If all checks pass, it sets a current_user variable with the user object.
	*/

	return func(c *gin.Context) {

		bearerToken := c.Request.Header.Get("Authorization")

		if !strings.Contains(bearerToken, "Bearer") {
			response := helper.APIresponse("Unauthorized", http.StatusUnauthorized, "error", nil, "Cannot hit this endpoint with no authentication.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		splitBearer := strings.Split(bearerToken, " ")
		var tokenString string = ""
		if len(splitBearer) == 2 {
			tokenString = splitBearer[1]
		}
		accessToken, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIresponse("Unauthorized", http.StatusUnauthorized, "error", nil, "Cannot hit this endpoint with no authentication.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		payload, ok := accessToken.Claims.(jwt.MapClaims)

		if !ok || !accessToken.Valid {
			response := helper.APIresponse("Unauthorized", http.StatusUnauthorized, "error", nil, gin.H{"message": "Do not have permissions to access this resources"})
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		userID := payload["_id"].(string)
		user, err := userService.FindByID(c.Request.Context(), userID)

		if err != nil {
			response := helper.APIresponse("Unauthorized", http.StatusUnauthorized, "error", nil, "Cannot hit this endpoint with no authentication.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return

		}

		if user.Email != payload["email"] {
			response := helper.APIresponse("Unauthorized", http.StatusUnauthorized, "error", nil, "Cannot hit this endpoint with no authentication.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		c.Set("current_user", user)
	}
}
