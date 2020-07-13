package crud

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
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

func (dr *databaseSchoolRepo) GetSchools(searchText string) []domain.School {
	var schools []domain.School
	s := strings.Split(searchText, "")
	pattern := "%" + strings.Join(s, "%") + "%"
	dr.database.DB.Where("name LIKE ?", pattern).Find(&schools)
	return schools
}

func (dr *databaseSchoolRepo) Create(school *domain.School) error {
	return dr.database.DB.Create(school).Error
}

func (dr *databaseSchoolRepo) Update(school *domain.School) error {
	if dr.database.DB.Model(school).Update(school).RowsAffected == 0 {
		return messages.New("user_not_found", "user not found")
	}
	return nil
}

func (dr *databaseSchoolRepo) Delete(id uint) error {
	// TODO usar el soft delete de GORM
	return nil
}
