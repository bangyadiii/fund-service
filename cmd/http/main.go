package main

import (
	"backend-crowdfunding/config"
	"backend-crowdfunding/database"
	_ "backend-crowdfunding/docs"
	"backend-crowdfunding/insfrastructure/cache"
	"backend-crowdfunding/insfrastructure/firebase"
	"backend-crowdfunding/sdk/id"
	"backend-crowdfunding/sdk/shutdown"
	"backend-crowdfunding/src/handler"
	"backend-crowdfunding/src/repository"
	"backend-crowdfunding/src/service"
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

// @title  Crowdfunding platform API
// @version 1.0
// @description Fund is a crowdfunding platform that enables clients to post their projects in search of funding. With Fund, clients can create campaigns and showcase their ideas to a community of potential investors. This project built with Go

// @host localhost:8000
// @BasePath /api/v1
// @schemes http
func main() {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	ctx := context.Background()

	// setup config
	env := config.New(".env")
	cleanup, err := run(ctx, env)
	defer cleanup()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	shutdown.Gracefully(ctx)
}

func run(ctx context.Context, env config.Config) (func(), error) {
	server, cleanup, err := buildServer(env)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// start the server
	go func() {
		server.Run()
	}()
	// return a function to close the server and database
	return func() {
		cleanup()
		gracefullyShutdown(ctx, server)
	}, nil
}

func buildServer(env config.Config) (*handler.Rest, func(), error) {
	// init database
	db, err := database.InitPostgresSQL(env)
	if err != nil {
		return nil, func() {}, err
	}

	//setup Redis
	redisCache := cache.Init(env)

	// setup id generator
	var IDGenerator = id.NewUlid()

	firebaseAuth := firebase.NewFirebase()

	// setup repository
	repo := repository.InitRepository(db, redisCache, IDGenerator, firebaseAuth)

	// setup service
	svc := service.InitService(env, repo)

	// init handler
	rest := handler.Init(svc, env)
	return rest, func() {
		database.CloseDB(db)
	}, nil
}

func gracefullyShutdown(ctx context.Context, rest *handler.Rest) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rest.Shutdown(shutdownCtx)
}
