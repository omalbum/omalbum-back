package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/domain"
	"github.com/omalbum/omalbum-back/internal/api/messages"
	"github.com/omalbum/omalbum-back/internal/api/services/crud"
	"github.com/omalbum/omalbum-back/internal/api/services/permissions"
	"github.com/omalbum/omalbum-back/internal/api/services/problems"
	"github.com/omalbum/omalbum-back/internal/api/utils/params"
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
	cache    domain.TeleOMACache
}

func NewProblemsController(database *db.Database, manager permissions.Manager, cache domain.TeleOMACache) ProblemsController {
	logger := crud.NewLogger(database)
	return &problemsController{
		database: database,
		manager:  manager,
		logger:   logger,
		cache:    cache,
	}
}

func (p problemsController) GetNextProblems(context *gin.Context) {
	ps, err := problems.NewService(p.database, p.cache).GetNextProblems()
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, ps)
}

func (p problemsController) GetCurrentProblems(context *gin.Context) {
	ps, err := problems.NewService(p.database, p.cache).GetCurrentProblems()
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, ps)

}

func (p problemsController) GetProblem(context *gin.Context) {
	problemId := params.GetProblemID(context)
	problem, err := problems.NewService(p.database, p.cache).GetProblemById(problemId)
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, problem)
}

func (p problemsController) GetAllProblems(context *gin.Context) {
	ps, err := problems.NewService(p.database, p.cache).GetAllProblems()
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, ps)
}
