package main

import (
	"backend-crowdfunding/campaign"
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_NAME"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("error : ", err.Error())
	}
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&campaign.Campaign{})
	db.AutoMigrate(&campaign.CampaignImage{})

	router := routes.GetRouter(db)

	appAddress := fmt.Sprintf("%s:%s", os.Getenv("APP_ADDRESS"), os.Getenv("PORT"))
	router.GinRouter.Run(appAddress)
}
