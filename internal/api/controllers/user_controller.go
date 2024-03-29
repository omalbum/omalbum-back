package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/domain"
	"github.com/omalbum/omalbum-back/internal/api/messages"
	"github.com/omalbum/omalbum-back/internal/api/services/crud"
	"github.com/omalbum/omalbum-back/internal/api/services/mailer"
	"github.com/omalbum/omalbum-back/internal/api/services/permissions"
	"github.com/omalbum/omalbum-back/internal/api/services/users"
	"github.com/omalbum/omalbum-back/internal/api/utils/params"
	"net/http"
)

type UserController interface {
	GetUser(context *gin.Context)
	PutUser(context *gin.Context)
	PutPassword(context *gin.Context)
	ResetPassword(context *gin.Context)
	GetAlbum(context *gin.Context)
	GetProblemAttemptsByUser(context *gin.Context)
	PostAnswer(context *gin.Context)
}

type userController struct {
	database *db.Database
	manager  permissions.Manager
	mailer   mailer.Mailer
	logger   crud.Logger
}

func NewUserController(database *db.Database, manager permissions.Manager, mailer mailer.Mailer) UserController {
	logger := crud.NewLogger(database)
	return &userController{
		database: database,
		manager:  manager,
		mailer:   mailer,
		logger:   logger,
	}
}

// Answers a /users request
// Sends the user personal data and the data for each role the user has.
// Permissions to see the user profile are verified here
func (u *userController) GetUser(context *gin.Context) {
	userID := params.GetUserID(context)

	if !u.manager.IsAdminOrSameUser(context, userID) { //todo esto podria ser un middleware
		context.JSON(http.StatusForbidden, gin.H{})
		return
	}

	userApp, err := users.NewService(u.database, u.mailer).GetByUserID(userID)
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}

	context.JSON(http.StatusOK, userApp)
}
func (u *userController) PutUser(context *gin.Context) {

	userID := params.GetUserID(context) // we want to update this user's profile

	if !u.manager.IsAdminOrSameUser(context, userID) { //todo esto podria ser un middleware
		context.JSON(http.StatusForbidden, gin.H{})
		return
	}
	var updatedProfile domain.RegistrationApp // same payload format as register user
	err := context.Bind(&updatedProfile)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	userService := users.NewService(u.database, u.mailer)
	user, err := userService.UpdateUserProfile(userID, updatedProfile)

	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, user)
}

func (u *userController) PutPassword(context *gin.Context) {
	userID := params.GetUserID(context) // we want to update this user's password

	if !u.manager.IsAdminOrSameUser(context, userID) { //todo esto podria ser un middleware
		context.JSON(http.StatusForbidden, gin.H{})
		return
	}
	var passwordChangePayload domain.PasswordChangeApp
	err := context.Bind(&passwordChangePayload)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	err = users.NewService(u.database, u.mailer).ChangePassword(userID, passwordChangePayload.OldPassword, passwordChangePayload.NewPassword)

	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, gin.H{})
}

func (u *userController) ResetPassword(context *gin.Context) {
	var emailWrappedApp domain.EmailWrappedApp
	err := context.Bind(&emailWrappedApp)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	email := emailWrappedApp.Email
	err = users.NewService(u.database, u.mailer).ResetPassword(email)

	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, gin.H{})
}

func (u *userController) GetAlbum(context *gin.Context) {
	userID := params.GetUserID(context)
	if !u.manager.IsAdminOrSameUser(context, userID) { //todo esto podria ser un middleware
		context.JSON(http.StatusForbidden, gin.H{})
		return
	}
	album, err := users.NewService(u.database, u.mailer).GetAlbum(userID)
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, *album)
}

func (u *userController) GetProblemAttemptsByUser(context *gin.Context) {
	userID := params.GetUserID(context)
	if !u.manager.IsAdminOrSameUser(context, userID) { //todo esto podria ser un middleware
		context.JSON(http.StatusForbidden, gin.H{})
		return
	}
	problemID := params.GetProblemID(context)
	problemStats, err := users.NewService(u.database, u.mailer).GetProblemAttemptsByUser(userID, problemID)
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, *problemStats)
}

func (u *userController) PostAnswer(context *gin.Context) {
	userID := params.GetCallerID(context)
	var problemAttemptApp domain.ProblemAttemptApp
	err := context.Bind(&problemAttemptApp)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	result, err := users.NewService(u.database, u.mailer).PostAnswer(userID, problemAttemptApp)
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, *result)
}
