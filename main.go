package main

import (
	"backend-crowdfunding/auth"
	"backend-crowdfunding/handler"
	"backend-crowdfunding/user"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load(".env")
	dsn := "root:root@tcp(127.0.0.1:3306)/crowd_startup?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("error : ", err.Error())
	}
	db.AutoMigrate(&user.User{})

	repository := user.NewRepository(db)
	userService := user.NewService(repository)
	authService := auth.NewService()
	userHandler := handler.NewUserHanlder(userService, authService)

	router := gin.Default()

	api := router.Group("/api/v1")
	authApi := api.Group("/auth")
	authApi.POST("/email-is-available", userHandler.CheckIsEmailAvailable)
	authApi.POST("/register", userHandler.RegisterUser)
	authApi.POST("/login", userHandler.Login)
	authApi.POST("/avatars", userHandler.UploadAvatar)

	router.Run()	
}
