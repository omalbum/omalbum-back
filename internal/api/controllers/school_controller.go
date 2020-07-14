package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/mailer"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/permissions"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/schools"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/params"
	"net/http"
)

type SchoolController interface {
	GetSchools(context *gin.Context)
}

type schoolController struct {
	database *db.Database
	manager  permissions.Manager
	mailer   mailer.Mailer
	logger   crud.Logger
}

func NewSchoolController(database *db.Database, manager permissions.Manager, mailer mailer.Mailer) SchoolController {
	logger := crud.NewLogger(database)
	return &schoolController{
		database: database,
		manager:  manager,
		mailer:   mailer,
		logger:   logger,
	}
}

func (u *schoolController) GetSchools(context *gin.Context) {
	searchText := params.GetSearchText(context)
	province := params.GetProvince(context)
	department := params.GetDepartment(context)
	schoolsApp, err := schools.NewService(u.database, u.mailer).GetSchools(searchText, province, department)
	if err != nil {
		context.JSON(messages.GetHttpCode(err), err)
		return
	}

	context.JSON(http.StatusOK, schoolsApp)
}
