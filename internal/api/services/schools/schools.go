package schools

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/mailer"
)

type Service interface {
	GetSchools(searchText string, province string, department string) ([]domain.School, error)
}

type service struct {
	database *db.Database
	mailer   mailer.Mailer
}

func NewService(database *db.Database, mailer mailer.Mailer) Service {
	return &service{
		database: database,
		mailer:   mailer,
	}
}

func (s *service) GetSchools(searchText string, province string, department string) ([]domain.School, error) {
	schools := crud.NewDatabaseSchoolRepo(s.database).GetSchools(searchText, province, department)
	return schools, nil
}