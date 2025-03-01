package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasen-devvv/mini-blog-backend/config"
	"github.com/jasen-devvv/mini-blog-backend/models"
)

// CommentInput struktur untuk input komentar
type CommentInput struct {
	Content string `json:"content" binding:"required"`
}

// GetComments mengambil semua komentar untuk artikel
func GetComments(c *gin.Context) {
	articleID := c.Param("id")

	var comments []models.Comment

	// Preload user untuk mendapatkan info penulis komentar
	if err := config.DB.Where("article_id = ?", articleID).Preload("User").Order("created_at asc").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	// Hapus password dari respons
	for i := range comments {
		comments[i].User.Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

// CreateComment membuat komentar baru
func CreateComment(c *gin.Context) {
	articleID := c.Param("id")

	// Dapatkan user_id dari context (dari middleware auth)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Periksa apakah artikel ada
	var article models.Article
	if err := config.DB.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// Bind input
	var input CommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buat komentar baru
	comment := models.Comment{
		Content:   input.Content,
		UserID:    userID.(uint),
		ArticleID: article.ID,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Load info user untuk respons
	config.DB.Preload("User").First(&comment, comment.ID)

	// Hapus password dari respons
	comment.User.Password = ""

	c.JSON(http.StatusCreated, gin.H{"data": comment})
}
