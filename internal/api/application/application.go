package application

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/miguelsotocarlos/teleoma/internal/api/clients/sendgrid"
	"github.com/miguelsotocarlos/teleoma/internal/api/config"
	"github.com/miguelsotocarlos/teleoma/internal/api/controllers"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/middlewares"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/cache"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/mailer"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/permissions"
	"github.com/spf13/afero"
	"time"
)

type Application struct {
	AuthMiddleware        middlewares.AuthMiddleware
	AdminMiddleware       middlewares.AdminMiddleware
	UserProblemMiddleware middlewares.UserProblemMiddleware
	RegisterController    controllers.RegisterController
	UserController        controllers.UserController
	AdminController       controllers.AdminController
	ProblemsController    controllers.ProblemsController
	//SchoolController      controllers.SchoolController
}

func BuildApplication(db *db.Database) *Application {
	manager := permissions.NewManager(db)

	var sendGridRestClient sendgrid.RestClient
	if config.ShouldSendMails() {
		sendGridRestClient = sendgrid.NewRestClient(resty.New(), config.GetSendGridApiKey())
	} else {
		sendGridRestClient = sendgrid.NewRestClientMock()
	}
	mail := mailer.New(sendGridRestClient, mailer.NewTemplateLoader(afero.NewOsFs()))
	teleOMACache := cache.NewTeleOMACache()

	authMiddleware := middlewares.NewAuthMiddleware(db)
	adminMiddleware := middlewares.NewAdminMiddleware(db, manager)
	userProblemMiddleware := middlewares.NewUserProblemMiddleware(db, manager, teleOMACache)
	registerController := controllers.NewRegisterController(db, mail)
	userController := controllers.NewUserController(db, manager, mail)
	adminController := controllers.NewAdminController(db, manager)
	problemsController := controllers.NewProblemsController(db, manager, teleOMACache)
	//schoolController := controllers.NewSchoolController(db)
	return &Application{
		AuthMiddleware:        authMiddleware,
		AdminMiddleware:       adminMiddleware,
		UserProblemMiddleware: userProblemMiddleware,
		RegisterController:    registerController,
		UserController:        userController,
		AdminController:       adminController,
		ProblemsController:    problemsController,
		//SchoolController:      schoolController,
	}
}

func Setup(db *db.Database) *gin.Engine {
	engine := gin.Default()

	if config.IsCorsEnabled() {
		engine.Use(corsHandler())
	}

	app := BuildApplication(db)
	MapURLs(app, engine)

	return engine
}

func corsHandler() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
