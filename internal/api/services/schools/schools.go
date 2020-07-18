package schools

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
)

type Service interface {
	GetSchools(searchText string, province string, department string) []domain.SchoolApp
}

type service struct {
	database *db.Database
}

func NewService(database *db.Database) Service {
	return &service{
		database: database,
	}
}

func (s *service) GetSchools(searchText string, province string, department string) []domain.SchoolApp {
	schools := crud.NewDatabaseSchoolRepo(s.database).GetSchools(searchText, province, department)
	schoolsApp := make([]domain.SchoolApp, len(schools))
	for i, school := range schools {
		schoolsApp[i] = schoolToSchoolApp(school)
	}
	return schoolsApp
}

func schoolToSchoolApp(school domain.School) domain.SchoolApp {
	return domain.SchoolApp{Name: school.Name, Province: school.Province, Department: school.Department, Location: school.Location}
}
