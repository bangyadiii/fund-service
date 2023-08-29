package handler

import (
	"backend-crowdfunding/config"
	"backend-crowdfunding/src/middleware"
	"backend-crowdfunding/src/service"
	"context"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var once = sync.Once{}

type rest struct {
	http    *fiber.App
	service *service.Service
	cfg     config.Config
}

func Init(s *service.Service, cfg config.Config) *rest {
	r := &rest{}
	once.Do(func() {
		r.http = fiber.New()
		r.cfg = cfg
		r.service = s

		r.RegisterMiddlewareAndRoutes()
	})

	return r
}

func (r *rest) RegisterMiddlewareAndRoutes() {
	// Global middleware
	r.http.Use(recover.New())
	r.http.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	//r.http.Use(swagger.New(r.cfg))

	// auth router group
	api := r.http.Group("api/v1")
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

func (r *rest) Run() {
	// set parent context
	c := context.Background()
	ctx, stop := signal.NotifyContext(c, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	/*
		Create context that listens for the interrupt signal from the OS.
		This will allow us to gracefully shutdown the server.
	*/

	var port = r.cfg.GetWithDefault("APP_PORT", ":8080")

	go func() {
		if err := r.http.Listen(port); err != nil && err != http.ErrServerClosed {
			log.Fatal("error while listening server,", err)
		}
	}()

	log.Printf("Server is running at %s", port)

	<-ctx.Done()
	stop()
	log.Println("Shutting down server")
	shutdownCtx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := r.http.ShutdownWithContext(shutdownCtx); err != nil {
		log.Fatalf("error while shutting down server: %v", err)
	}

	log.Println("The server has been shutdown")
}
