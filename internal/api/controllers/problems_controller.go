package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/permissions"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/problems"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/params"
	"net/http"
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
	ps, err := problems.NewService(p.database).GetNextProblems()
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, ps)
}

func (p problemsController) GetCurrentProblems(context *gin.Context) {
	ps, err := problems.NewService(p.database).GetCurrentProblems()
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, ps)

}

func (p problemsController) GetProblem(context *gin.Context) {
	problemId := params.GetProblemID(context)
	problem, err := problems.NewService(p.database).GetProblemById(problemId)
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, problem)
}

func (p problemsController) GetAllProblems(context *gin.Context) {
	panic("implement me") //TODO
}
