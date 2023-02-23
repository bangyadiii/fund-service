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
	return func(c *gin.Context) {
		// TODO mengecek request header dengan key Authorization
		//  TODO mengecek, apakah value Authorization terdapat "Bearer" value
		// TODO split string dengan spasi (" ")
		// TODO mengambil array index ke 2
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
		userID := uint(payload["_id"].(float64))
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
