package router

import (
	"time"

	"github.com/gin-contrib/cors"
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

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5500",
			"http://127.0.0.1:5500",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
