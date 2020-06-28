package admin

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"time"
)

type Service interface {
	CreateProblem(poserId uint, newProblem domain.NewProblemApp) (*domain.Problem, error)
	UpdateProblem(problemId uint, updatedProblem domain.NewProblemApp) (*domain.Problem, error)
	DeleteProblem(problemId uint) error
}

type service struct {
	database *db.Database
}

func (s *service) DeleteProblem(problemId uint) error {
	var err = crud.NewDatabaseProblemRepo(s.database).Delete(problemId)
	if err != nil {
		return messages.NewNotFound("problem_not_found", "problem not found")
	}
	return nil
}

func NewService(database *db.Database) Service {
	return &service{
		database: database,
	}
}

func (s *service) CreateProblem(poserId uint, newProblem domain.NewProblemApp) (*domain.Problem, error) {
	problemRepo := crud.NewDatabaseProblemRepo(s.database)
	err := newProblem.Validate()
	if err != nil {
		return nil, messages.NewValidation(err)
	}
	problem := problemAppToProblem(newProblem)
	problem.PoserId = poserId
	err = problemRepo.Create(&problem)
	// TODO falta guardar las tags
	return &problem, err
}

func (s *service) UpdateProblem(problemId uint, updatedProblem domain.NewProblemApp) (*domain.Problem, error) {
	problemRepo := crud.NewDatabaseProblemRepo(s.database)
	err := updatedProblem.Validate()
	if err != nil {
		return nil, messages.NewValidation(err)
	}
	problem := problemAppToProblem(updatedProblem)
	problem.ID = problemId
	err = problemRepo.Update(&problem)
	// todo UPDATE TAGS
	return &problem, err
}

func problemAppToProblem(newProblem domain.NewProblemApp) domain.Problem {
	return domain.Problem{
		OmaforosPostId:   newProblem.OmaforosPostId,
		DateUploaded:     time.Now(),
		DateContestStart: newProblem.ReleaseDate,
		DateContestEnd:   newProblem.Deadline,
		Statement:        newProblem.Statement,
		Answer:           newProblem.Answer,
		Annotations:      newProblem.Annotations,
		IsDraft:          newProblem.IsDraft,
		Hint:             newProblem.Hint,
		OfficialSolution: newProblem.OfficialSolution,
	}
}
