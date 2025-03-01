package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jasen-devvv/mini-blog-backend/controllers"
)

func SetupAuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}
}
