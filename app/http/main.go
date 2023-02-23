package main

import (
	"backend-crowdfunding/config"
	"backend-crowdfunding/database/migrations"
	"backend-crowdfunding/src/handler"
	"backend-crowdfunding/src/repository"
	"backend-crowdfunding/src/service"
	"backend-crowdfunding/src/util/id"
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
	m.RunMigration()

	// setup id generator
	var IDGenerator = id.NewUlid()

	// setup repository
	repo := repository.InitRepository(db, IDGenerator)

	// setup service
	service := service.InitService(configuration, repo)

	// init handler
	rest := handler.Init(service, configuration)

	rest.Run()
}
