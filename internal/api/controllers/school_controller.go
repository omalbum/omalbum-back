package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/services/schools"
	"github.com/omalbum/omalbum-back/internal/api/utils/params"
	"net/http"
)

type SchoolController interface {
	GetSchools(context *gin.Context)
}

type schoolController struct {
	database *db.Database
}

func NewSchoolController(database *db.Database) SchoolController {
	return &schoolController{
		database: database,
	}
}

func (u *schoolController) GetSchools(context *gin.Context) {
	searchText := params.GetSearchText(context)
	province := params.GetProvince(context)
	department := params.GetDepartment(context)
	schoolsApp := schools.NewService(u.database).GetSchools(searchText, province, department)
	context.JSON(http.StatusOK, schoolsApp)
}
