package database

import (
	"backend-crowdfunding/config"
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func InitPostgresSQL(env config.Config) (*DB, error) {
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
		return &DB{}, err
	}

	postgresDB, err := db.DB()

	if err != nil {
		return &DB{}, err
	}
	maxIdle, _ := strconv.Atoi(env.Get("DB_MAX_IDLE"))
	OpenCon, _ := strconv.Atoi(env.Get("DB_MAX_CONN"))
	MaxIdleTime, _ := strconv.Atoi(env.Get("DB_MAX_IDLE_TIME_IN_MINUTES"))
	maxLifetime, _ := strconv.Atoi(env.Get("DB_MAX_LIFETIME_IN_MINUTES"))

	postgresDB.SetMaxIdleConns(maxIdle)
	postgresDB.SetMaxOpenConns(OpenCon)
	postgresDB.SetConnMaxIdleTime(time.Duration(MaxIdleTime) * time.Minute)
	postgresDB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Minute)

	return &DB{DB: db}, nil
}

func CloseDB(db *DB) {
	sqlDB, _ := db.DB.DB()

	_ = sqlDB.Close()
	log.Println("database has been closed.")
}
