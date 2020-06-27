package application

import (
	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MapURLs(app *Application, router *gin.Engine) {

	v1 := router.Group("/api/v1")
	{
		// Ping
		v1.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		// Authentication endpoints
		// Login
		v1.POST("/auth/login", app.AuthMiddleware.LoginHandler)
		// Refresh token
		v1.POST("/auth/refresh_token", app.AuthMiddleware.RefreshHandler)

		// Register endpoints
		v1.POST("/register", app.RegisterController.Register)
	}

	jobs := router.Group("/jobs")
	{
		// Resource to return the JSON data
		jobs.GET("/json", JobJson)
	}

	router.NoRoute()
}

// returns a map[string]interface{} that can be marshalled as JSON
func JobJson(c *gin.Context) {
	c.JSON(200, jobrunner.StatusJson())
}
