package problems

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"time"
)

type Service interface {
	GetProblemById(problemId uint) (domain.ProblemApp, error)
	GetCurrentProblems() (domain.CurrentProblemsApp, error)
	GetNextProblems() (domain.NextProblemsApp, error)
	GetAllProblems() (domain.AllProblemsApp, error)
}

type service struct {
	database *db.Database
	cache    domain.TeleOMACache
}

func (s *service) GetAllProblems() (domain.AllProblemsApp, error) {
	var res = s.cache.Get(domain.AllProblemsCacheKey)
	if res != nil {
		return res.(domain.AllProblemsApp), nil
	}
	//TODO Optimizacion: traer solamente las tags necesarias con un inner join
	problemsDatabase := crud.NewDatabaseProblemRepo(s.database).GetAllProblems()
	tagRepo := crud.NewDatabaseProblemTagRepo(s.database)
	problems := make([]domain.ProblemApp, len(problemsDatabase))
	position := make(map[uint]int)
	for i, p := range problemsDatabase {
		position[p.ID] = i
	}
	for i, prob := range problemsDatabase {
		problems[i] = problemToProblemApp(prob)
	}
	var problemTags = tagRepo.GetAllTags()
	for _, tag := range problemTags {
		problemId := tag.ProblemId
		if i, ok := position[problemId]; ok {
			problems[i].Tags = append(problems[i].Tags, tag.Tag)
		}
	}
	allProblems := domain.AllProblemsApp{Problems: problems}
	s.cache.SetWithTTL(domain.AllProblemsCacheKey, allProblems, time.Second)
	return allProblems, nil
}

func (s *service) GetNextProblems() (domain.NextProblemsApp, error) {
	var res = s.cache.Get(domain.NextProblemsCacheKey)
	if res != nil {
		return res.(domain.NextProblemsApp), nil
	}
	problemsDatabase := crud.NewDatabaseProblemRepo(s.database).GetNextProblems()
	problems := make([]domain.ProblemNextApp, len(problemsDatabase))

	for i, prob := range problemsDatabase {
		problems[i] = problemToProblemNextApp(prob)
	}
	nextProblems := domain.NextProblemsApp{NextProblems: problems}
	s.cache.SetWithTTL(domain.NextProblemsCacheKey, nextProblems, time.Second)
	return nextProblems, nil
}

func (s *service) GetCurrentProblems() (domain.CurrentProblemsApp, error) {
	var res = s.cache.Get(domain.CurrentProblemsCacheKey)
	if res != nil {
		return res.(domain.CurrentProblemsApp), nil
	}
	//TODO Optimizacion: traer tags con un inner join para hacer menos queries a la DB
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
	currentProblems := domain.CurrentProblemsApp{CurrentProblems: problems}
	s.cache.SetWithTTL(domain.CurrentProblemsCacheKey, currentProblems, time.Second)
	return currentProblems, nil
}

func NewService(database *db.Database, cache domain.TeleOMACache) Service {
	return &service{
		database: database,
		cache:    cache,
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
