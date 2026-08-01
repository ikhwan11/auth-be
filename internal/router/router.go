package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ikhwan11/auth-be/internal/container"
)

func Setup(app *container.Container) *gin.Engine {
	r := gin.Default()

	// Swagger
	r.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "Auth API is running",
		})
	})

	auth := r.Group("/auth")
	{
		auth.POST(
			"/check-employee",
			app.AuthHandler.CheckEmployee,
		)

		auth.POST(
			"/register",
			app.AuthHandler.Register,
		)

		auth.POST(
			"/login",
			app.AuthHandler.Login,
		)

		auth.POST(
			"/refresh",
			app.AuthHandler.RefreshToken,
		)

		auth.POST("/logout", app.AuthHandler.Logout)
	}

	return r
}
