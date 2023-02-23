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

	"github.com/gin-gonic/gin"
)

var once = sync.Once{}

type rest struct {
	http    *gin.Engine
	service service.Service
	cfg     config.Config
}

func Init(s *service.Service, cfg config.Config) *rest {
	r := &rest{}
	once.Do(func() {
		gin.SetMode(gin.ReleaseMode) // TODO: Move to config later

		r.http = gin.New()
		r.cfg = cfg

		r.RegisterMiddlewareAndRoutes()
	})

	return r
}

func (r *rest) RegisterMiddlewareAndRoutes() {
	// Global middleware

	// auth router group
	api := r.http.Group("api/v1")
	authApi := r.http.Group("/auth")
	authApi.POST("/email-is-available", r.CheckIsEmailAvailable)
	authApi.POST("/register", r.RegisterUser)
	authApi.POST("/login", r.Login)
	authApi.POST("/avatars", middleware.VerifyToken(r.service.User, r.service.Auth), r.UploadAvatar)

	// campaign router group
	campaignRoute := api.Group("/campaigns")

	campaignImages := campaignRoute.Group("/images")
	campaignImages.POST("/", middleware.VerifyToken(r.service.User, r.service.Auth), r.UploadCampaignImage)

	campaignRoute.PUT("/:id", middleware.VerifyToken(r.service.User, r.service.Auth), r.UpdateCampaign)
	campaignRoute.GET("/:id", r.GetCampaignByID)
	campaignRoute.GET("/", r.GetCampaigns)
	campaignRoute.POST("/", middleware.VerifyToken(r.service.User, r.service.Auth), r.CreateNewCampaign)

	trxRoutes := api.Group("/model.Transactions")
	trxRoutes.GET("/", r.GetAllTransactionsByCampaignID)
	trxRoutes.POST("/", middleware.VerifyToken(r.service.User, r.service.Auth), r.CreateTransaction)
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

	port := ":8080"
	if r.cfg.Get("APP_PORT") != "" {
		port = ":" + r.cfg.Get("APP_PORT")
	}

	server := &http.Server{
		Addr:              port,
		Handler:           r.http,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("error while listening server,", err)
		}
	}()

	log.Printf("Server is running at %s", port)
	<-ctx.Done()
	stop()
	log.Println("Shutting down server")
	shutdownCtx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("error while shutting down server: %v", err)
	}

	log.Println("The server has been shutdown")
}
