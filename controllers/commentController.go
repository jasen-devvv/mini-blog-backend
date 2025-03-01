package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasen-devvv/mini-blog-backend/config"
	"github.com/jasen-devvv/mini-blog-backend/models"
)

// CommentInput defines the structure for comment creation requests
type CommentInput struct {
	Content string `json:"content" binding:"required"`
}

// GetComments retrieves all comments for a specific article, ordered by creation time.
// Comments include user information with passwords removed for security.
// Returns a JSON response with the comments or an error message.
func GetComments(c *gin.Context) {
	articleID := c.Param("id")
	var comments []models.Comment

	// Preload user data to get comment author information
	if err := config.DB.Where("article_id = ?", articleID).Preload("User").Order("created_at asc").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	// Remove password from user data for security
	for i := range comments {
		comments[i].User.Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

// CreateComment adds a new comment to an article.
// Requires authentication, as it uses the user_id from the context (set by auth middleware).
// Validates that the referenced article exists before creating the comment.
// Returns a JSON response with the created comment or an appropriate error message.
func CreateComment(c *gin.Context) {
	articleID := c.Param("id")

	// Get user_id from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Verify the article exists before adding a comment
	var article models.Article
	if err := config.DB.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// Validate input
	var input CommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create new comment
	comment := models.Comment{
		Content:   input.Content,
		UserID:    userID.(uint),
		ArticleID: article.ID,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Load user info for response
	config.DB.Preload("User").First(&comment, comment.ID)

	// Remove password from response for security
	comment.User.Password = ""

	c.JSON(http.StatusCreated, gin.H{"data": comment})
}
