package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jasen-devvv/mini-blog-backend/controllers"
	"github.com/jasen-devvv/mini-blog-backend/middleware"
)

// SetupArticleRoutes sets up the article-related routes for the application.
//
// Available routes:
//   - GET    /api/articles       -> Fetch all articles
//   - GET    /api/articles/:id   -> Fetch a specific article by ID
//   - POST   /api/articles       -> Create a new article (requires authentication)
//   - PUT    /api/articles/:id   -> Update an existing article by ID (requires authentication)
//   - DELETE /api/articles/:id   -> Delete an article by ID (requires authentication)
//
// Routes that modify data (POST, PUT, DELETE) are protected by authentication middleware.
func SetupArticleRoutes(router *gin.Engine) {
	articles := router.Group("/api/articles")
	{
		articles.GET("", controllers.GetAllArticles)
		articles.GET("/:id", controllers.GetArticle)

		articles.Use(middleware.AuthMiddleware())
		{
			articles.POST("", controllers.CreateArticle)
			articles.PUT("/:id", controllers.UpdateArticle)
			articles.DELETE("/:id", controllers.DeleteArticle)
		}
	}
}
