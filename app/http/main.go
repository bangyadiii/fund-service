package main

import (
	"backend-crowdfunding/config"
	"backend-crowdfunding/database/migrations"
	"backend-crowdfunding/insfrastructure/firebase"
	"backend-crowdfunding/sdk/id"
	"backend-crowdfunding/src/handler"
	"backend-crowdfunding/src/repository"
	"backend-crowdfunding/src/service"
	"log"
)

func main() {
	// setup config
	configuration := config.New(".env")

	// init database
	db, err := config.InitPostgreSQL(configuration)
	if err != nil {
		log.Fatalf("error when connect to db, %v", err)
	}

	m := migrations.Migration{DB: db}

	// run migration
	err = m.RunMigration()
	if err != nil {
		log.Fatalf("Migration error, %v", err)
		return
	}

	//setup Redis
	//redisIns := cache.Init(configuration)

	// setup id generator
	var IDGenerator = id.NewUlid()

	firebaseAuth := firebase.NewFirebase()

	// setup repository
	repo := repository.InitRepository(db, IDGenerator, firebaseAuth)

	// setup service
	svc := service.InitService(configuration, repo)

	// init handler
	rest := handler.Init(svc, configuration)

	rest.Run()
}
