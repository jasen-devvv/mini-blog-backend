package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jasen-devvv/mini-blog-backend/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_URL")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	DB.AutoMigrate(&models.User{}, &models.Article{}, &models.Comment{})

	fmt.Println("Database connected successfully")
}
