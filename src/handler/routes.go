package handler

import (
	"backend-crowdfunding/src/middleware"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (r *Rest) RegisterMiddlewareAndRoutes() {
	// Global middleware
	r.Http.Use(recover.New())
	r.Http.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	//r.Http.Use(swagger.New(r.cfg))

	// auth router group
	api := r.Http.Group("api/v1")
	authApi := api.Group("/auth")
	authApi.Post("/email-is-available", r.CheckIsEmailAvailable)
	authApi.Post("/register", r.RegisterUser)
	authApi.Post("/login", r.Login)
	authApi.Post("/login/google", r.LoginWithGoogle)
	authApi.Post("/avatars", middleware.VerifyToken(r.service.User, r.service.Auth), r.UploadAvatar)

	// campaign router group
	campaignRoute := api.Group("/campaigns")

	campaignImages := campaignRoute.Group("/images")

	campaignRoute.Get("/:id", r.GetCampaignByID)
	campaignRoute.Get("/", r.GetCampaigns)
	// authenticated require routes
	campaignImages.Post("/", middleware.VerifyToken(r.service.User, r.service.Auth), r.UploadCampaignImage)
	campaignRoute.Put("/:id", middleware.VerifyToken(r.service.User, r.service.Auth), r.UpdateCampaign)
	campaignRoute.Post("/", middleware.VerifyToken(r.service.User, r.service.Auth), r.CreateNewCampaign)

	trxRoutes := api.Group("/transactions")
	trxRoutes.Get("/", r.GetAllTransactionsByCampaignID)
	trxRoutes.Post("/", middleware.VerifyToken(r.service.User, r.service.Auth), r.CreateTransaction)
}
