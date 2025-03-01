package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasen-devvv/mini-blog-backend/config"
	"github.com/jasen-devvv/mini-blog-backend/models"
)

// ArticleInput defines the structure for article creation and update requests
type ArticleInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// GetAllArticles retrieves all articles from the database, ordered by creation date (newest first)
// and includes the associated user information with the password field removed.
// Returns a JSON response with the articles or an error message.
func GetAllArticles(ctx *gin.Context) {
	var articles []models.Article

	if err := config.DB.Preload("User").Order("created_at desc").Find(&articles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
		return
	}

	// Remove password from user data for security
	for i := range articles {
		articles[i].User.Password = ""
	}

	ctx.JSON(http.StatusOK, gin.H{"data": articles})
}

// GetArticle retrieves a single article by its ID, including the associated user information.
// The password field is removed from the user data for security.
// Returns a JSON response with the article or a "not found" error.
func GetArticle(ctx *gin.Context) {
	id := ctx.Param("id")

	var article models.Article

	if err := config.DB.Preload("User").First(&article, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	article.User.Password = ""

	ctx.JSON(http.StatusOK, gin.H{"data": article})
}

// CreateArticle creates a new article in the database.
// Requires authentication, as it uses the user_id from the context (set by auth middleware).
// Returns a JSON response with the created article (including user info) or an error message.
func CreateArticle(c *gin.Context) {
	var input ArticleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user_id from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Create new article
	article := models.Article{
		Title:   input.Title,
		Content: input.Content,
		UserID:  userID.(uint),
	}

	if err := config.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	// Load user info for the response
	config.DB.Preload("User").First(&article, article.ID)

	// Remove password from response for security
	article.User.Password = ""

	c.JSON(http.StatusCreated, gin.H{"data": article})
}

// UpdateArticle updates an existing article in the database.
// Requires authentication and verifies that the user is the owner of the article.
// Returns a JSON response with the updated article or an appropriate error message.
func UpdateArticle(c *gin.Context) {
	id := c.Param("id")

	// Get user_id from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if article exists
	var article models.Article
	if err := config.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// Check if user is the owner of the article
	if article.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this article"})
		return
	}

	// Bind input
	var input ArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update article
	if err := config.DB.Model(&article).Updates(models.Article{
		Title:   input.Title,
		Content: input.Content,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	// Load user info for response
	config.DB.Preload("User").First(&article, article.ID)

	// Remove password from response for security
	article.User.Password = ""

	c.JSON(http.StatusOK, gin.H{"data": article})
}

// DeleteArticle removes an article from the database.
// Requires authentication and verifies that the user is the owner of the article.
// Returns a success message or an appropriate error message.
func DeleteArticle(c *gin.Context) {
	id := c.Param("id")

	// Get user_id from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if article exists
	var article models.Article
	if err := config.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// Check if user is the owner of the article
	if article.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this article"})
		return
	}

	// Delete article
	if err := config.DB.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Article deleted successfully"})
}
