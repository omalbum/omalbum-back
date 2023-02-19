package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/domain"
	"github.com/omalbum/omalbum-back/internal/api/messages"
	"github.com/omalbum/omalbum-back/internal/api/services/crud"
	"github.com/omalbum/omalbum-back/internal/api/services/mailer"
	"github.com/omalbum/omalbum-back/internal/api/services/register"
	"net/http"
)

type RegisterController interface {
	Register(context *gin.Context)
}

type registerController struct {
	registerService register.Service
	logger          crud.Logger
}

func NewRegisterController(database *db.Database, mailer mailer.Mailer) RegisterController {
	return &registerController{
		registerService: register.NewService(database, mailer),
		logger:          crud.NewLogger(database),
	}
}

func (r *registerController) Register(context *gin.Context) {
	var registrationApp domain.RegistrationApp
	err := context.Bind(&registrationApp)
	if err != nil {
		_ = r.logger.LogAnonymousAction("cannot_bind", http.StatusBadRequest, context.Request.Method, context.Request.URL.String())
		context.JSON(http.StatusBadRequest, err)
		return
	}
	u, err := r.registerService.CreateUser(registrationApp)
	if err != nil {
		_ = r.logger.LogAnonymousAction(err.(messages.Message).Code, messages.GetHttpCode(err), context.Request.Method, context.Request.URL.String())
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	_ = r.logger.LogUserAction(u.UserId, "user_registered", http.StatusCreated, context.Request.Method, context.Request.URL.String())
	context.JSON(http.StatusCreated, u)
}
