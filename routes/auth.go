package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jasen-devvv/mini-blog-backend/controllers"
)

// SetupAuthRoutes sets up authentication-related routes for the application.
//
// Available routes:
//   - POST /api/auth/register -> Register a new user
//   - POST /api/auth/login    -> Authenticate and log in a user
func SetupAuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}
}
