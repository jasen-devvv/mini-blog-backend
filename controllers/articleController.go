package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasen-devvv/mini-blog-backend/config"
	"github.com/jasen-devvv/mini-blog-backend/models"
)

type ArticleInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func GetAllArticles(ctx *gin.Context) {
	var articles []models.Article

	if err := config.DB.Preload("User").Order("created_at desc").Find(&articles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
		return
	}

	for i := range articles {
		articles[i].User.Password = ""
	}

	ctx.JSON(http.StatusOK, gin.H{"data": articles})
}

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

// CreateArticle membuat artikel baru
func CreateArticle(c *gin.Context) {
	var input ArticleInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Dapatkan user_id dari context (dari middleware auth)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Buat artikel baru
	article := models.Article{
		Title:   input.Title,
		Content: input.Content,
		UserID:  userID.(uint),
	}

	if err := config.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	// Load info user untuk respons
	config.DB.Preload("User").First(&article, article.ID)
	
	// Hapus password dari respons
	article.User.Password = ""

	c.JSON(http.StatusCreated, gin.H{"data": article})
}

// UpdateArticle memperbarui artikel
func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	
	// Dapatkan user_id dari context (dari middleware auth)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Periksa apakah artikel ada
	var article models.Article
	if err := config.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// Periksa apakah user adalah pemilik artikel
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

	// Update artikel
	if err := config.DB.Model(&article).Updates(models.Article{
		Title:   input.Title,
		Content: input.Content,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	// Load info user untuk respons
	config.DB.Preload("User").First(&article, article.ID)
	
	// Hapus password dari respons
	article.User.Password = ""

	c.JSON(http.StatusOK, gin.H{"data": article})
}

// DeleteArticle menghapus artikel
func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	
	// Dapatkan user_id dari context (dari middleware auth)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Periksa apakah artikel ada
	var article models.Article
	if err := config.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// Periksa apakah user adalah pemilik artikel
	if article.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this article"})
		return
	}

	// Hapus artikel
	if err := config.DB.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Article deleted successfully"})
}