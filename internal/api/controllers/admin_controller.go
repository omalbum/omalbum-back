package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/admin"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/permissions"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/params"
	"net/http"
)

type AdminController interface {
	GetProblem(context *gin.Context)
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
		r.logger.LogAnonymousAction("problem creation failed: validation error", http.StatusBadRequest, context.Request.Method, context.Request.URL.String())
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
	_ = context.Bind(&updatedProblem)
	problemId := params.GetProblemID(context)
	_, err := admin.NewService(a.database).UpdateProblem(problemId, updatedProblem)
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}
	context.JSON(http.StatusOK, gin.H{})
}

func (a *adminController) DeleteProblem(context *gin.Context) {
	problemId := params.GetProblemID(context)
	err := admin.NewService(a.database).DeleteProblem(problemId)
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
	}
	context.JSON(http.StatusOK, gin.H{})
}
