package crud

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
)

type databaseExpandedUserProblemAttemptRepo struct {
	database *db.Database
}

func (d databaseExpandedUserProblemAttemptRepo) GetByUserId(userId uint) []domain.ExpandedUserProblemAttempt {
	if userId == 0 {
		return nil //needed to avoid an empty condition in Where
	}
	var attempts []domain.ExpandedUserProblemAttempt
	d.database.DB.Where("user_id = ?", userId).Find(&attempts)
	return attempts
}

func (d databaseExpandedUserProblemAttemptRepo) GetByUserAndProblemId(userId uint, problemId uint, isContestProblem bool) []domain.ExpandedUserProblemAttempt {
	if userId == 0 {
		return nil //needed to avoid an empty condition in Where
	}
	var attempts []domain.ExpandedUserProblemAttempt
	duringContestPossibilities := "(1, 0)"
	if isContestProblem {
		duringContestPossibilities = "(0)"
	}
	d.database.DB.Where("user_id = ? AND problem_id = ? AND during_contest IN ?", userId, problemId, duringContestPossibilities).Find(&attempts)
	return attempts
}

func NewExpandedUserProblemAttemptRepo(database *db.Database) domain.ExpandedUserProblemAttemptRepo {
	return &databaseExpandedUserProblemAttemptRepo{
		database: database,
	}
}
