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
	logger   crud.Logger
}

func NewRegisterController(database *db.Database, mailer mailer.Mailer) RegisterController {
	return &registerController{
		database: database,
		mailer:   mailer,
		logger:   crud.NewLogger(database),
	}
}

func (r *registerController) Register(context *gin.Context) {
	var registrationApp domain.RegistrationApp
	_ = context.Bind(&registrationApp)

	var err = registrationApp.Validate()

	if err != nil {
		r.logger.LogAnonymousAction("registration failed: validation error", http.StatusBadRequest, context.Request.Method, context.Request.URL.String())
		context.JSON(http.StatusBadRequest, messages.NewValidation(err))
		return
	}

	// Create User
	user, err := createUser(r.database, &registrationApp)
	if err != nil {
		r.logger.LogAnonymousAction("registration failed: user already exists", http.StatusConflict, context.Request.Method, context.Request.URL.String())
		context.JSON(http.StatusConflict, messages.New("user_already_exists", "user already exists"))
		return
	}

	// Send the mail in a non-blocking way
	// comento porque por ahora no usamos esto
	// registrationJob := mailer.NewRegistrationJob(r.mailer, registrationApp.Email, registrationApp.Name)
	// jobrunner.Now(registrationJob)
	r.logger.LogUserAction(user.ID, "user registered", http.StatusOK, context.Request.Method, context.Request.URL.String())
	context.JSON(http.StatusCreated, gin.H{"user_id": user.ID})
}

func createUser(database *db.Database, registrationApp *domain.RegistrationApp) (*domain.User, error) {
	userRepo := crud.NewDatabaseUserRepo(database)
	birthDate, _ := time.Parse("2006-01-02", registrationApp.BirthDate)
	user := domain.User{
		UserName:         strings.ToLower(registrationApp.UserName),
		HashedPassword:   crypto.HashAndSalt(registrationApp.Password),
		RegistrationDate: time.Now(),
		LastActiveDate:   time.Now(),
		Name:             registrationApp.Name,
		LastName:         registrationApp.LastName,
		BirthDate:        birthDate,
		Email:            strings.ToLower(registrationApp.Email),
	}
	err := userRepo.Create(&user)

	return &user, err
}
