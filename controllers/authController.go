package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jasen-devvv/mini-blog-backend/config"
	"github.com/jasen-devvv/mini-blog-backend/models"
	"golang.org/x/crypto/bcrypt"
)

// RegisterInput defines the structure for user registration request
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput defines the structure for user login request
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register creates a new user account in the database.
// It validates the input, hashes the password, and returns the created user.
// Password validation ensures it's at least 6 characters long.
func Register(ctx *gin.Context) {
	var input RegisterInput

	// Validate input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password for secure storage
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user with hashed password
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	// Save user to database, handle potential duplicate email/username
	if err := config.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user. Email or username may already exist."})
		return
	}

	// Hide password in response for security
	user.Password = ""
	ctx.JSON(http.StatusCreated, gin.H{"data": user})
}

// Login authenticates a user and generates a JWT token.
// It checks email and password, and returns a token and user info on success.
// The token expires after one week.
func Login(ctx *gin.Context) {
	var input LoginInput

	// Validate input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verify password against stored hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token with user ID and expiration
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // Token expires in 1 week
	})

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Hide password in response for security
	user.Password = ""
	ctx.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  user,
	})
}
