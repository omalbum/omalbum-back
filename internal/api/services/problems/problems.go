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
	GetNextProblems() ([]domain.ProblemNextApp, error)
}

type service struct {
	database *db.Database
}

func (s *service) GetNextProblems() ([]domain.ProblemNextApp, error) {
	problemsDatabase := crud.NewDatabaseProblemRepo(s.database).GetNextProblems()
	problems := make([]domain.ProblemNextApp, len(problemsDatabase))

	for i, prob := range problemsDatabase {
		problems[i] = problemToProblemNextApp(prob)
	}
	return problems, nil
}

func (s *service) GetCurrentProblems() ([]domain.ProblemApp, error) {
	//TODO Optimizacion: traer roles con un inner join para hacer menos queries a la DB
	// eso es lo costoso de esta funcion.
	problemsDatabase := crud.NewDatabaseProblemRepo(s.database).GetCurrentProblems()
	tagRepo := crud.NewDatabaseProblemTagRepo(s.database)
	problems := make([]domain.ProblemApp, len(problemsDatabase))

	for i, prob := range problemsDatabase {
		problems[i] = problemToProblemApp(prob)
		var problemTags = tagRepo.GetByProblemId(prob.ID)
		var tags = make([]string, len(problemTags))
		for j, problemTag := range problemTags {
			tags[j] = problemTag.Tag
		}
		problems[i].Tags = tags
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
		ProblemId:      problem.ID,
		OmaforosPostId: problem.OmaforosPostId,
		ReleaseDate:    problem.DateContestStart,
		Deadline:       problem.DateContestEnd,
		Statement:      problem.Statement,
	}
}

func problemToProblemNextApp(problem domain.Problem) domain.ProblemNextApp {
	return domain.ProblemNextApp{
		ProblemId:   problem.ID,
		ReleaseDate: problem.DateContestStart,
		Deadline:    problem.DateContestEnd,
	}
}
