package problems

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
)

type Service interface {
	GetProblemById(problemId uint) (domain.ProblemApp, error)
	GetCurrentProblems() ([]domain.ProblemApp, error)
}

type service struct {
	database *db.Database
}

func (s *service) GetCurrentProblems() ([]domain.ProblemApp, error) {
	problemsDatabase := crud.NewDatabaseProblemRepo(s.database).GetCurrentProblems()
	problems := make([]domain.ProblemApp, len(problemsDatabase))
	for i, prob := range problemsDatabase {
		problems[i] = problemToProblemApp(prob)
	}
	return problems, nil
}

func NewService(database *db.Database) Service {
	return &service{
		database: database,
	}
}

func (s *service) GetProblemById(problemId uint) (domain.ProblemApp, error) {
	problem := crud.NewDatabaseProblemRepo(s.database).GetById(problemId)
	if problem == nil {
		return domain.ProblemApp{}, messages.NewNotFound("problem_not_found", "problem not found")
	}
	problemApp := problemToProblemApp(*problem)
	problemApp.ProblemId = problemId
	problemTags := crud.NewDatabaseProblemTagRepo(s.database).GetByProblemId(problemId)
	var tags = make([]string, len(problemTags))
	for i, problemTag := range problemTags {
		tags[i] = problemTag.Tag
	}
	problemApp.Tags = tags
	return problemApp, nil
}

func problemToProblemApp(problem domain.Problem) domain.ProblemApp {
	return domain.ProblemApp{
		OmaforosPostId: problem.OmaforosPostId,
		ReleaseDate:    problem.DateContestStart,
		Deadline:       problem.DateContestEnd,
		Statement:      problem.Statement,
	}
}
