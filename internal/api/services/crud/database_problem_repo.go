package crud

import (
	"github.com/jinzhu/gorm"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"time"
)

type databaseProblemRepo struct {
	database *db.Database
}

func NewDatabaseProblemRepo(database *db.Database) domain.ProblemRepo {
	return &databaseProblemRepo{
		database: database,
	}
}

func (dr *databaseProblemRepo) GetById(ID uint) *domain.Problem {
	if ID == 0 {
		return nil //needed to avoid an empty condition in Where
	}
	var problem domain.Problem
	if dr.database.DB.Where(&domain.Problem{Model: gorm.Model{ID: ID}}).First(&problem).RecordNotFound() {
		return nil
	}

	return &problem
}

func (dr *databaseProblemRepo) Create(problem *domain.Problem) error {
	return dr.database.DB.Create(problem).Error
}

func (dr *databaseProblemRepo) Update(problem *domain.Problem) error {
	if problem.ID == 0 {
		return messages.New("problem_id_must_be_nonzero", "problem id must be nonzero")
	}
	if dr.database.DB.Model(problem).Update(problem).RowsAffected == 0 {
		return messages.New("user_not_found", "user not found")
	}
	return nil
}

func (dr *databaseProblemRepo) Delete(problemId uint) error {
	if problemId == 0 { // esto es clave para no borrar la tabla entera accidentalmente
		return messages.New("problem_id_must_be_nonzero", "problem id must be nonzero")
	}
	return dr.database.DB.Where(&domain.Problem{Model: gorm.Model{ID: problemId}}).Delete(&domain.Problem{Model: gorm.Model{ID: problemId}}).Error
}

func (dr *databaseProblemRepo) GetNextProblems() []domain.Problem {
	t := time.Now()
	var problems []domain.Problem
	dr.database.DB.Where("(? < problems.date_contest_start) AND (NOT problems.is_draft)", t).Find(&problems)
	return problems

}

func (dr *databaseProblemRepo) GetCurrentProblems() []domain.Problem {
	t := time.Now()
	var problems []domain.Problem
	dr.database.DB.Where("(problems.date_contest_start < ?) AND (? < problems.date_contest_end) AND (NOT problems.is_draft)", t, t).Find(&problems)
	return problems
}

func (dr *databaseProblemRepo) GetAllProblems() []domain.Problem {
	t := time.Now()
	var problems []domain.Problem
	dr.database.DB.Where("(problems.date_contest_start < ? ) AND  (NOT problems.is_draft)", t).Find(&problems)
	return problems
}
