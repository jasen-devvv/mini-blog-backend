package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jasen-devvv/mini-blog-backend/controllers"
	"github.com/jasen-devvv/mini-blog-backend/middleware"
)

func SetupArticleRoutes(router *gin.Engine) {
	articles := router.Group("/api/articles")
	{
		articles.GET("", controllers.GetAllArticles)
		articles.GET("/:id", controllers.GetArticle)

		// Protected routes
		articles.Use(middleware.AuthMiddleware())
		{
			articles.POST("", controllers.CreateArticle)
			articles.PUT("/:id", controllers.UpdateArticle)
			articles.DELETE("/:id", controllers.DeleteArticle)
		}
	}
}
