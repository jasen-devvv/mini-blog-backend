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

// DB is the global database connection instance used throughout the application
var DB *gorm.DB

// ConnectDatabase initializes the database connection.
//
// It performs the following steps:
// - Loads environment variables from .env file
// - Reads the database connection URL from the environment variable DB_URL
// - Connects to the PostgreSQL database using GORM
// - Runs automatic migrations for User, Article, and Comment models
//
// If any step fails, the application will log an error and terminate.
func ConnectDatabase() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get database connection string from environment
	dsn := os.Getenv("DB_URL")

	// Initialize database connection
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// AutoMigrate ensures tables exist and updates schema if necessary
	DB.AutoMigrate(&models.User{}, &models.Article{}, &models.Comment{})

	fmt.Println("Database connected successfully")
}
