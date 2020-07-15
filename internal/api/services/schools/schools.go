package schools

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
)

type Service interface {
	GetSchools(searchText string, province string, department string) []domain.School
}

type service struct {
	database *db.Database
}

func NewService(database *db.Database) Service {
	return &service{
		database: database,
	}
}

func (s *service) GetSchools(searchText string, province string, department string) []domain.School {
	schools := crud.NewDatabaseSchoolRepo(s.database).GetSchools(searchText, province, department)
	return schools
}
