package main

import (
	"backend-crowdfunding/routes"
	"backend-crowdfunding/user"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	    log.Fatal("Error loading .env file")
	}
	
	dsn := "root:root@tcp(127.0.0.1:3306)/crowd_startup?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("error : ", err.Error())
	}
	db.AutoMigrate(&user.User{})

	router := routes.GetRouter(db)

	appAddress := fmt.Sprintf("%s:%s", os.Getenv("APP_ADDRESS"), os.Getenv("PORT")) 
	router.GinRouter.Run(appAddress)
}

