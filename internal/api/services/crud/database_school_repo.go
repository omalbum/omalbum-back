package crud

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"strings"
)

type databaseSchoolRepo struct {
	database *db.Database
}

func NewDatabaseSchoolRepo(database *db.Database) domain.SchoolRepo {
	return &databaseSchoolRepo{
		database: database,
	}
}

func (dr *databaseSchoolRepo) GetSchools(searchText string, province string, department string) []domain.School {
	var schools []domain.School
	matchExact := true
	var pattern string
	if matchExact {
		pattern = "%" + searchText + "%"
	} else {
		s := strings.Split(searchText, "")
		pattern = "%" + strings.Join(s, "%") + "%"
	}
	dr.database.DB.Where("name LIKE ? AND province = ? AND department = ?", pattern, province, department).Find(&schools)
	return schools
}

func (dr *databaseSchoolRepo) Create(school *domain.School) error {
	return dr.database.DB.Create(school).Error
}

