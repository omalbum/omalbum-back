package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/permissions"
)

type AdminController interface {
	PostProblem(context *gin.Context)
	PutProblem(context *gin.Context)
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
