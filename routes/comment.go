package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jasen-devvv/mini-blog-backend/controllers"
	"github.com/jasen-devvv/mini-blog-backend/middleware"
)

// SetupCommentRoutes sets up comment-related routes for the application.
//
// Available routes:
//   - GET  /api/articles/:id/comments  -> Fetch all comments for an article
//   - POST /api/articles/:id/comments  -> Add a new comment to an article (requires authentication)
//
// The POST route is protected by authentication middleware.
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
