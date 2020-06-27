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

		//Problems endpoints

		//Users endpoints
		v1.GET("/users/:user_id/profile", app.AuthMiddleware.MiddlewareFunc(), app.UserController.GetUser)
		v1.PUT("/users/:user_id/profile", app.AuthMiddleware.MiddlewareFunc(), app.UserController.PutUser)
		v1.PUT("/users/:user_id/password", app.AuthMiddleware.MiddlewareFunc(), app.UserController.PutPassword)
		v1.POST("/users/password/reset", app.UserController.ResetPassword) //TODO este no esta funcional sin configurar el envio de emails!

		//Admin endpoints
		//v1.POST("/admin/problem", app.AuthMiddleware.MiddlewareFunc(), app.AdminMiddleware.AdminCheck, app.AdminController.PostProblem)
		v1.POST("/admin/problem", app.AuthMiddleware.MiddlewareFunc(), app.AdminMiddleware.AdminCheck, app.AdminController.PostProblem)

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
