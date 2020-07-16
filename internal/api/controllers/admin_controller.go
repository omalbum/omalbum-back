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
	GetAllProblems(context *gin.Context)
	GetAllProblemsStats(context *gin.Context)
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

func (a *adminController) GetAllProblemsStats(context *gin.Context) {
	problems := admin.NewService(a.database).GetAllProblems()
	var problemAdminStats = make([]domain.ProblemAdminStatsApp, len(problems.Problems))
	for i, problem := range problems.Problems {
		allAttempts := admin.NewService(a.database).GetAllProblemAttempts(problem.ProblemId)
		problemAdminStats[i].ProblemId = problem.ProblemId
		problemAdminStats[i].NumberInSeries = problem.NumberInSeries
		problemAdminStats[i].IsCurrentProblem = problem.IsCurrentProblemAdminApp()
		problemAdminStats[i].Tags = problem.Tags
		setUsers := make(map[uint]bool)
		for _, attempt := range allAttempts {
			problemAdminStats[i].Attempts++
			if attempt.IsCorrect {
				problemAdminStats[i].SolvedCount++
				setUsers[attempt.UserId] = true
				if attempt.DuringContest {
					problemAdminStats[i].SolvedDuringContestCount++
				}
			}
		}
		problemAdminStats[i].SolvedDistinctCount = len(setUsers)
	}
	context.JSON(http.StatusOK, problemAdminStats)
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
	_ = context.Bind(&updatedProblem)
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
