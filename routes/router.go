package routes

import (
	"backend-crowdfunding/auth"
	"backend-crowdfunding/handler"
	"backend-crowdfunding/middleware"
	"backend-crowdfunding/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type router struct{
	GinRouter *gin.Engine
	db *gorm.DB
}


func GetRouter(db *gorm.DB) router{
	router := newRouter(gin.Default(), db)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	userHandler := handler.NewUserHanlder(userService, authService)

	// campaignRepository := campaign.NewRepository(db)
	// campaignService := campaign.NewService(campaignRepository)
	// campaignHandler := handler.NewUserHanlder(campaignService, authService)



	api := router.GinRouter.Group("/api/v1")
	authApi := api.Group("/auth")
	authApi.POST("/email-is-available", userHandler.CheckIsEmailAvailable)
	authApi.POST("/register", userHandler.RegisterUser)
	authApi.POST("/login", userHandler.Login)
	authApi.POST("/avatars", middleware.VerifyToken(userService, authService), userHandler.UploadAvatar)


	return *router
}

func newRouter(ginEngine *gin.Engine, db *gorm.DB) *router{
	return &router{ginEngine,db}
}