package config

import (
	"fmt"
	"strconv"
	"time"

	"backend-crowdfunding/database"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgreSQL(configuration Config) (*database.DB, error) {
	var err error

	connStr := fmt.Sprintf("user=%s password=%s port=%s dbname=%s sslmode=%s",
		configuration.Get("DB_USER"),
		configuration.Get("DB_PASSWORD"),
		configuration.Get("DB_PORT"),
		configuration.Get("DB_NAME"),
		configuration.Get("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(connStr))

	if err != nil {
		return &database.DB{}, err
	}

	postgreDB, err := db.DB()

	if err != nil {
		return &database.DB{}, err
	}
	maxIdle, _ := strconv.Atoi(configuration.Get("DB_MAX_IDLE"))
	OpenCon, _ := strconv.Atoi(configuration.Get("DB_MAX_CONN"))
	Max_idle_time, _ := strconv.Atoi(configuration.Get("DB_MAX_IDLE_TIME_IN_MINUTES"))
	max_lifetime, _ := strconv.Atoi(configuration.Get("DB_MAX_LIFETIME_IN_MINUTES"))

	postgreDB.SetMaxIdleConns(maxIdle)
	postgreDB.SetMaxOpenConns(OpenCon)
	postgreDB.SetConnMaxIdleTime(time.Duration(Max_idle_time) * time.Minute)
	postgreDB.SetConnMaxLifetime(time.Duration(max_lifetime) * time.Minute)

	return &database.DB{DB: db}, nil
}
