package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jasen-devvv/mini-blog-backend/controllers"
	"github.com/jasen-devvv/mini-blog-backend/middleware"
)

func SetupCommentRoutes(router *gin.Engine) {
	// Public routes
	router.GET("/api/articles/:id/comments", controllers.GetComments)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/articles/:id/comments", controllers.CreateComment)
	}
}
