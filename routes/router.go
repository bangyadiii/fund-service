package routes

import (
	"backend-crowdfunding/src/handler"
	"backend-crowdfunding/src/middleware"
	"backend-crowdfunding/src/repository"
	"backend-crowdfunding/src/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type router struct {
	GinRouter *gin.Engine
	db        *gorm.DB
}

func GetRouter(db *gorm.DB) router {
	router := newRouter(gin.Default(), db)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService()
	userHandler := handler.NewUserHanlder(userService, authService)

	campaignRepository := repository.NewCampaignRepository(db)
	campaignService := service.NewCampaignService(campaignRepository)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	trxRepo := repository.NewTransactionRepository(db)
	trxService := service.NewTransactionService(trxRepo)
	trxHandler := handler.NewTransactionHandler(trxService)

	api := router.GinRouter.Group("/api/v1")

	// auth router group
	authApi := api.Group("/auth")
	authApi.POST("/email-is-available", userHandler.CheckIsEmailAvailable)
	authApi.POST("/register", userHandler.RegisterUser)
	authApi.POST("/login", userHandler.Login)
	authApi.POST("/avatars", middleware.VerifyToken(userService, authService), userHandler.UploadAvatar)

	// campaign router group
	campaignRoute := api.Group("/campaigns")

	campaignImages := campaignRoute.Group("/images")
	campaignImages.POST("/", middleware.VerifyToken(userService, authService), campaignHandler.UploadCampaignImage)

	campaignRoute.PUT("/:id", middleware.VerifyToken(userService, authService), campaignHandler.UpdateCampaign)
	campaignRoute.GET("/:id", campaignHandler.GetCampaignByID)
	campaignRoute.GET("/", campaignHandler.GetCampaigns)
	campaignRoute.POST("/", middleware.VerifyToken(userService, authService), campaignHandler.CreateNewCampaign)

	trxRoutes := api.Group("/model.Transactions")
	trxRoutes.GET("/", trxHandler.GetAllTransactionsByCampaignID)
	trxRoutes.POST("/", middleware.VerifyToken(userService, authService), trxHandler.CreateTransaction)

	return *router
}

func newRouter(ginEngine *gin.Engine, db *gorm.DB) *router {
	return &router{ginEngine, db}
}
