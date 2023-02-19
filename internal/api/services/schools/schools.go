package schools

import (
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/domain"
	"github.com/omalbum/omalbum-back/internal/api/services/crud"
)

type Service interface {
	GetSchools(searchText string, province string, department string) []string
}

type service struct {
	database *db.Database
}

func NewService(database *db.Database) Service {
	return &service{
		database: database,
	}
}

func (s *service) GetSchools(searchText string, province string, department string) []string {
	schools := crud.NewDatabaseSchoolRepo(s.database).GetSchools(searchText, province, department)
	schoolsApp := make([]string, len(schools))
	for i, school := range schools {
		schoolsApp[i] = schoolToSchoolApp(school)
	}
	return schoolsApp
}

func schoolToSchoolApp(school domain.School) string {
	return school.Name
	//return domain.SchoolApp{Name: school.Name, Province: school.Province, Department: school.Department, Location: school.Location}
}
