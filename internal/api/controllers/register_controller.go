package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/mailer"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/crypto"
	"net/http"
	"strings"
	"time"
)

type RegisterController interface {
	Register(context *gin.Context)
}

type registerController struct {
	database *db.Database
	mailer   mailer.Mailer
}

func NewRegisterController(database *db.Database, mailer mailer.Mailer) RegisterController {
	return &registerController{
		database: database,
		mailer:   mailer,
	}
}

func (r *registerController) Register(context *gin.Context) {
	var registrationApp domain.RegistrationApp
	_ = context.Bind(&registrationApp)

	var err = registrationApp.Validate()
	logger := crud.NewLogger(r.database)

	if err != nil {
		registrationApp.Password = "******"
		logger.LogAnonymousAction("registration failed: validation error", registrationApp)
		context.JSON(http.StatusBadRequest, messages.NewValidation(err))
		return
	}

	// Create User
	user, err := createUser(r.database, &registrationApp)
	if err != nil {
		registrationApp.Password = "******"
		logger.LogAnonymousAction("registration failed: user already exists", registrationApp)
		context.JSON(http.StatusConflict, messages.New("user_already_exists", "user already exists"))
		return
	}

	// log user action
	logger.LogUserActionAtTime(user.ID, "user registered", "", user.RegistrationDate)

	// Send the mail in a non-blocking way
	// comento porque por ahora no usamos esto
	// registrationJob := mailer.NewRegistrationJob(r.mailer, registrationApp.Email, registrationApp.Name)
	// jobrunner.Now(registrationJob)

	context.JSON(http.StatusCreated, gin.H{"user_id": user.ID})
}

func createUser(database *db.Database, registrationApp *domain.RegistrationApp) (*domain.User, error) {
	userRepo := crud.NewDatabaseUserRepo(database)
	user := domain.User{
		UserName:         strings.ToLower(registrationApp.UserName),
		HashedPassword:   crypto.HashAndSalt(registrationApp.Password),
		RegistrationDate: time.Now(),
		LastActiveDate:   time.Now(),
		Name:             registrationApp.Name,
		LastName:         registrationApp.LastName,
		Cellphone:        registrationApp.Cellphone,
		Email:            strings.ToLower(registrationApp.Email),
	}
	err := userRepo.Create(&user)

	return &user, err
}
