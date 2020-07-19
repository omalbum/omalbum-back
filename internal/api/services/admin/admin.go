package admin

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
)

type Service interface {
	GetAllProblems() domain.AllProblemsAdminApp
	GetProblem(problemId uint) (*domain.ProblemAdminApp, error)
	CreateProblem(poserId uint, newProblem domain.ProblemAdminApp) (*domain.Problem, error)
	UpdateProblem(problemId uint, updatedProblem domain.ProblemAdminApp) (*domain.ProblemAdminApp, error)
	DeleteProblem(problemId uint) error
	GetStats() []domain.ProblemAdminStatsApp
}

type service struct {
	database *db.Database
}

func (s *service) GetStats() []domain.ProblemAdminStatsApp {
	problems := s.GetAllProblems()
	var problemAdminStats = make([]domain.ProblemAdminStatsApp, len(problems.Problems))
	for i, problem := range problems.Problems {
		allAttempts := crud.NewExpandedUserProblemAttemptRepo(s.database).GetAllByProblem(problem.ProblemId)
		problemAdminStats[i].ProblemId = problem.ProblemId
		problemAdminStats[i].NumberInSeries = problem.NumberInSeries
		problemAdminStats[i].IsCurrentProblem = problem.IsCurrentProblemAdminApp()
		problemAdminStats[i].Tags = problem.Tags
		setUsers := make(map[uint]bool)
		for _, attempt := range allAttempts {
			problemAdminStats[i].Attempts++
			if attempt.IsCorrect {
				problemAdminStats[i].SolvedCount++
				setUsers[attempt.UserId] = true
				if attempt.DuringContest {
					problemAdminStats[i].SolvedDuringContestCount++
				}
			}
		}
		problemAdminStats[i].SolvedDistinctCount = len(setUsers)
	}
	return problemAdminStats
}

func (s *service) GetAllProblems() domain.AllProblemsAdminApp {
	problemsDatabase := crud.NewDatabaseProblemRepo(s.database).GetAllProblemsForAdmin()
	problems := make([]domain.ProblemAdminApp, len(problemsDatabase))
	position := make(map[uint]int)
	for i, p := range problemsDatabase {
		position[p.ID] = i
	}
	for i, prob := range problemsDatabase {
		problems[i] = problemToProblemAdminApp(prob)
	}
	tagRepo := crud.NewDatabaseProblemTagRepo(s.database)
	var problemTags = tagRepo.GetAllTags()
	for _, tag := range problemTags {
		problemId := tag.ProblemId
		if i, ok := position[problemId]; ok {
			problems[i].Tags = append(problems[i].Tags, tag.Tag)
		}
	}
	allProblemsAdmin := domain.AllProblemsAdminApp{Problems: problems}
	return allProblemsAdmin
}

func (s *service) GetProblem(problemId uint) (*domain.ProblemAdminApp, error) {
	problem := crud.NewDatabaseProblemRepo(s.database).GetById(problemId)
	if problem == nil {
		return nil, messages.NewNotFound("problem_not_found", "problem not found")
	}
	problemAdminApp := problemToProblemAdminApp(*problem)
	problemAdminApp.ProblemId = problemId
	problemTags := crud.NewDatabaseProblemTagRepo(s.database).GetByProblemId(problemId)
	var tags = make([]string, len(problemTags))
	for i, problemTag := range problemTags {
		tags[i] = problemTag.Tag
	}
	problemAdminApp.Tags = tags
	return &problemAdminApp, nil
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

func (s *service) UpdateProblem(problemId uint, updatedProblem domain.ProblemAdminApp) (*domain.ProblemAdminApp, error) {
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
	return s.GetProblem(problem.ID)
}

func problemAdminAppToProblem(newProblem domain.ProblemAdminApp) domain.Problem {
	return domain.Problem{
		OmaforosPostId:   newProblem.OmaforosPostId,
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
		ProblemId:        problem.ID,
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
