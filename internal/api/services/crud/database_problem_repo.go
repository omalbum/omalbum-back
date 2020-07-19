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

func (dr *databaseProblemRepo) GetNextNumberInSeries(series string) uint {
	var problem domain.Problem
	if dr.database.DB.Order("number_in_series desc").Where("(is_draft=0) and (series = ?)", series).First(&problem).RecordNotFound() {
		return 1
	}
	return problem.NumberInSeries + 1
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
	if problem.Series == "" {
		return messages.New("series_must_be_nonempty", "series must be nonempty")
	}
	if problem.IsDraft {
		problem.NumberInSeries = 0
	} else {
		problem.NumberInSeries = dr.GetNextNumberInSeries(problem.Series)
	}
	return dr.database.DB.Create(problem).Error
}

func (dr *databaseProblemRepo) Update(problem *domain.Problem) error {
	if problem.ID == 0 {
		return messages.New("problem_id_must_be_nonzero", "problem id must be nonzero")
	}
	var problemOld = dr.GetById(problem.ID)
	if problemOld == nil {
		return messages.New("problem_not_found", "problem not found")
	}
	if problem.IsDraft && !problemOld.IsDraft {
		return messages.New("cannot_convert_to_draft", "cannot convert problem to draft")
	}
	if problemOld.IsDraft && !problem.IsDraft {
		problem.NumberInSeries = dr.GetNextNumberInSeries(problem.Series)
	}
	if !problemOld.IsDraft {
		problem.Series = problemOld.Series // this should not be changed!
		problem.NumberInSeries = problemOld.NumberInSeries
	}
	if dr.database.DB.Model(problem).Update(problem).RowsAffected == 0 {
		return messages.New("unknown_error_updating_problem", "unknown error updating problem")
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

func (dr *databaseProblemRepo) GetAllProblemsForAdmin() []domain.Problem {
	var problems []domain.Problem
	dr.database.DB.Find(&problems)
	return problems
}
