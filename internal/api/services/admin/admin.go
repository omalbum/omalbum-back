package admin

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"time"
)

type Service interface {
	GetProblem(problemId uint) (domain.ProblemAdminApp, error)
	CreateProblem(poserId uint, newProblem domain.ProblemAdminApp) (*domain.Problem, error)
	UpdateProblem(problemId uint, updatedProblem domain.ProblemAdminApp) (*domain.Problem, error)
	DeleteProblem(problemId uint) error
}

type service struct {
	database *db.Database
}

func (s *service) GetProblem(problemId uint) (domain.ProblemAdminApp, error) {
	problem := crud.NewDatabaseProblemRepo(s.database).GetById(problemId)
	if problem == nil {
		return domain.ProblemAdminApp{}, messages.NewNotFound("problem_not_found", "problem not found")
	}
	problemAdminApp := problemToProblemAdminApp(*problem)
	problemAdminApp.ProblemId = problemId
	problemTags := crud.NewDatabaseProblemTagRepo(s.database).GetByProblemId(problemId)
	var tags = make([]string, len(problemTags))
	for i, problemTag := range problemTags {
		tags[i] = problemTag.Tag
	}
	problemAdminApp.Tags = tags
	return problemAdminApp, nil
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

func (s *service) CreateProblem(poserId uint, newProblem domain.ProblemAdminApp) (*domain.Problem, error) {
	problemRepo := crud.NewDatabaseProblemRepo(s.database)
	err := newProblem.Validate()
	if err != nil {
		return nil, messages.NewValidation(err)
	}
	problem := problemAdminAppToProblem(newProblem)
	problem.PoserId = poserId
	err = problemRepo.Create(&problem)
	if err != nil {
		return nil, messages.NewBadRequest("error", "error")
	}
	problemTagRepo := crud.NewDatabaseProblemTagRepo(s.database)
	problemTagRepo.CreateByProblemIdAndTags(problem.ID, newProblem.Tags)
	return &problem, nil
}

func (s *service) UpdateProblem(problemId uint, updatedProblem domain.ProblemAdminApp) (*domain.Problem, error) {
	problemRepo := crud.NewDatabaseProblemRepo(s.database)
	err := updatedProblem.Validate()
	if err != nil {
		return nil, messages.NewValidation(err)
	}
	problem := problemAdminAppToProblem(updatedProblem)
	problem.ID = problemId
	err = problemRepo.Update(&problem)
	if err != nil {
		return nil, messages.NewNotFound("problem_not_found", "problem not found")
	}
	problemTagRepo := crud.NewDatabaseProblemTagRepo(s.database)
	problemTagRepo.DeleteAllTagsByProblemId(problemId)
	problemTagRepo.CreateByProblemIdAndTags(problemId, updatedProblem.Tags)
	return &problem, nil
}

func problemAdminAppToProblem(newProblem domain.ProblemAdminApp) domain.Problem {
	return domain.Problem{
		OmaforosPostId:   newProblem.OmaforosPostId,
		DateUploaded:     time.Now(),
		DateContestStart: newProblem.ReleaseDate,
		DateContestEnd:   newProblem.Deadline,
		Series:           newProblem.Series,
		Statement:        newProblem.Statement,
		Answer:           newProblem.Answer,
		Annotations:      newProblem.Annotations,
		IsDraft:          newProblem.IsDraft,
		Hint:             newProblem.Hint,
		OfficialSolution: newProblem.OfficialSolution,
	}
}

func problemToProblemAdminApp(problem domain.Problem) domain.ProblemAdminApp {
	return domain.ProblemAdminApp{
		OmaforosPostId:   problem.OmaforosPostId,
		ReleaseDate:      problem.DateContestStart,
		Deadline:         problem.DateContestEnd,
		Statement:        problem.Statement,
		Series:           problem.Series,
		NumberInSeries:   problem.NumberInSeries,
		Answer:           problem.Answer,
		Annotations:      problem.Annotations,
		IsDraft:          problem.IsDraft,
		Hint:             problem.Hint,
		OfficialSolution: problem.OfficialSolution,
	}
}
