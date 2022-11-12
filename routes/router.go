package routes

import (
	"backend-crowdfunding/auth"
	"backend-crowdfunding/campaign"
	"backend-crowdfunding/handler"
	"backend-crowdfunding/middleware"
	"backend-crowdfunding/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type router struct {
	GinRouter *gin.Engine
	db        *gorm.DB
}

func GetRouter(db *gorm.DB) router {
	router := newRouter(gin.Default(), db)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	userHandler := handler.NewUserHanlder(userService, authService)

	campaignRepository := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepository)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	api := router.GinRouter.Group("/api/v1")

	// auth router group
	authApi := api.Group("/auth")
	authApi.POST("/email-is-available", userHandler.CheckIsEmailAvailable)
	authApi.POST("/register", userHandler.RegisterUser)
	authApi.POST("/login", userHandler.Login)
	authApi.POST("/avatars", middleware.VerifyToken(userService, authService), userHandler.UploadAvatar)

	// campaign router group
	campaignRoute := api.Group("/campaigns")
	campaignRoute.GET("/", campaignHandler.GetCampaigns)
	campaignRoute.GET("/:id", campaignHandler.GetCampaignByID)
	campaignRoute.POST("/", middleware.VerifyToken(userService, authService), campaignHandler.CreateNewCampaign)

	return *router
}

func newRouter(ginEngine *gin.Engine, db *gorm.DB) *router {
	return &router{ginEngine, db}
}
