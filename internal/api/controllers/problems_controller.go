package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/permissions"
)

type ProblemsController interface {
	GetNextProblems(context *gin.Context)
	GetCurrentProblems(context *gin.Context)
	GetProblem(context *gin.Context)
	GetAllProblems(context *gin.Context)
}

type problemsController struct {
	database *db.Database
	manager  permissions.Manager
	logger   crud.Logger
}

func NewProblemsController(database *db.Database, manager permissions.Manager) ProblemsController {
	logger := crud.NewLogger(database)
	return &problemsController{
		database: database,
		manager:  manager,
		logger:   logger,
	}
}

func (p problemsController) GetNextProblems(context *gin.Context) {
	panic("implement me") //TODO
}

func (p problemsController) GetCurrentProblems(context *gin.Context) {
	panic("implement me") //TODO
}

func (p problemsController) GetProblem(context *gin.Context) {
	panic("implement me") //TODO
}

func (p problemsController) GetAllProblems(context *gin.Context) {
	panic("implement me") //TODO
}
