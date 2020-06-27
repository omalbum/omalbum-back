package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/permissions"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/params"
	"net/http"
)

type AdminController interface {
	PostProblem(context *gin.Context)
	PutProblem(context *gin.Context)
	DeleteProblem(context *gin.Context)
}

type adminController struct {
	database *db.Database
	manager  permissions.Manager
	logger   crud.Logger
}

func NewAdminController(database *db.Database, manager permissions.Manager) AdminController {
	logger := crud.NewLogger(database)
	return &adminController{
		database: database,
		manager:  manager,
		logger:   logger,
	}
}

func (a *adminController) PostProblem(context *gin.Context) {
	panic("implement me")
}

func (a *adminController) PutProblem(context *gin.Context) {
	panic("implement me")
}

func (a *adminController) DeleteProblem(context *gin.Context) {
	problemId := params.GetProblemID(context)
	var err = crud.NewDatabaseProblemRepo(a.database).Delete(problemId)
	if err != nil {
		context.JSON(http.StatusNotFound, messages.NewValidation(err))
		return
	}
	context.JSON(http.StatusOK, gin.H{})
}
