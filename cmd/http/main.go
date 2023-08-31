package main

import (
	"backend-crowdfunding/config"
	"backend-crowdfunding/database/migrations"
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
	db, err := config.InitPostgresSQL(env)
	if err != nil {
		return nil, func() {}, nil
	}

	m := migrations.Migration{DB: db}

	// run migration
	err = m.RunMigration()
	if err != nil {
		return nil, func() {
			config.CloseDB(db)
		}, nil
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
		config.CloseDB(db)
	}, nil
}

func gracefullyShutdown(ctx context.Context, rest *handler.Rest) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rest.Shutdown(shutdownCtx)
}
