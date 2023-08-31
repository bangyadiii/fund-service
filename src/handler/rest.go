package handler

import (
	"backend-crowdfunding/config"
	"backend-crowdfunding/src/service"
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"sync"
)

var once = sync.Once{}

type Rest struct {
	Http    *fiber.App
	service *service.Service
	cfg     config.Config
}

func Init(s *service.Service, cfg config.Config) *Rest {
	r := &Rest{}
	once.Do(func() {
		r.Http = fiber.New()
		r.cfg = cfg
		r.service = s

		r.RegisterMiddlewareAndRoutes()
	})

	return r
}

func (r *Rest) Run() {
	var port = r.cfg.GetWithDefault("APP_PORT", "8080")

	if err := r.Http.Listen("0.0.0.0:" + port); err != nil {
		log.Fatal("error while listening server,", err)
	}
}

func (r Rest) Shutdown(ctx context.Context) {
	if err := r.Http.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("error while shutting down server: %v", err)
	}

	log.Println("The server has been shutdown")
}
