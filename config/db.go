package config

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"backend-crowdfunding/database"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgresSQL(env Config) (*database.DB, error) {
	var err error

	connStr := fmt.Sprintf("user=%s password=%s port=%s dbname=%s sslmode=%s",
		env.GetOrPanic("DB_USER"),
		env.GetOrPanic("DB_PASSWORD"),
		env.GetOrPanic("DB_PORT"),
		env.GetOrPanic("DB_NAME"),
		env.GetOrPanic("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(connStr))

	if err != nil {
		return &database.DB{}, err
	}

	postgresDB, err := db.DB()

	if err != nil {
		return &database.DB{}, err
	}
	maxIdle, _ := strconv.Atoi(env.Get("DB_MAX_IDLE"))
	OpenCon, _ := strconv.Atoi(env.Get("DB_MAX_CONN"))
	MaxIdleTime, _ := strconv.Atoi(env.Get("DB_MAX_IDLE_TIME_IN_MINUTES"))
	maxLifetime, _ := strconv.Atoi(env.Get("DB_MAX_LIFETIME_IN_MINUTES"))

	postgresDB.SetMaxIdleConns(maxIdle)
	postgresDB.SetMaxOpenConns(OpenCon)
	postgresDB.SetConnMaxIdleTime(time.Duration(MaxIdleTime) * time.Minute)
	postgresDB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Minute)

	return &database.DB{DB: db}, nil
}

func CloseDB(db *database.DB) {
	sqlDB, _ := db.DB.DB()

	_ = sqlDB.Close()
	log.Println("database has been closed.")
}
