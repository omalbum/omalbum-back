package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/domain"
	"github.com/omalbum/omalbum-back/internal/api/messages"
	"github.com/omalbum/omalbum-back/internal/api/services/admin"
	"github.com/omalbum/omalbum-back/internal/api/services/crud"
	"github.com/omalbum/omalbum-back/internal/api/services/permissions"
	"github.com/omalbum/omalbum-back/internal/api/utils/params"
	"net/http"
)

type AdminController interface {
	GetProblem(context *gin.Context)
	PostProblem(context *gin.Context)
	PutProblem(context *gin.Context)
	DeleteProblem(context *gin.Context)
	GetAllProblems(context *gin.Context)
	GetStats(context *gin.Context)
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

func (a *adminController) GetAllProblems(context *gin.Context) {
	problems := admin.NewService(a.database).GetAllProblems()
	context.JSON(http.StatusOK, problems)
}

func (a *adminController) GetStats(context *gin.Context) {
	context.JSON(http.StatusOK, admin.NewService(a.database).GetStats())
}

func (a *adminController) GetProblem(context *gin.Context) {
	problemId := params.GetProblemID(context)
	problem, err := admin.NewService(a.database).GetProblem(problemId)
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, problem)
}

func (a *adminController) PostProblem(context *gin.Context) {
	var newProblem domain.ProblemAdminApp
	var err = context.Bind(&newProblem)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	userId := params.GetCallerID(context)
	problem, err := admin.NewService(a.database).CreateProblem(userId, newProblem)
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, gin.H{"problem_id": problem.ID})
}

func (a *adminController) PutProblem(context *gin.Context) {
	var updatedProblem domain.ProblemAdminApp
	err := context.Bind(&updatedProblem)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	problemId := params.GetProblemID(context)
	updatedProblemObject, err := admin.NewService(a.database).UpdateProblem(problemId, updatedProblem)
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, updatedProblemObject)
}

func (a *adminController) DeleteProblem(context *gin.Context) {
	problemId := params.GetProblemID(context)
	err := admin.NewService(a.database).DeleteProblem(problemId)
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
	}
	context.JSON(http.StatusOK, gin.H{})
}
