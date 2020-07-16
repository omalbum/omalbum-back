package crud

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
)

type databaseUserProblemAttemptRepo struct {
	database *db.Database
}

func NewDatabaseUserProblemAttemptRepo(database *db.Database) domain.UserProblemAttemptRepo {
	return &databaseUserProblemAttemptRepo{
		database: database,
	}
}

func (dr *databaseUserProblemAttemptRepo) Create(attempt *domain.UserProblemAttempt) error {
	return dr.database.DB.Create(attempt).Error
}

func (dr *databaseUserProblemAttemptRepo) GetByProblemIdAndUserId(problemId uint, userId uint) []domain.UserProblemAttempt {
	if problemId == 0 {
		return nil //needed to avoid an empty condition in Where
	}
	var attempts []domain.UserProblemAttempt
	dr.database.DB.Where("problem_id = ? AND user_id = ?", problemId, userId).Find(&attempts)
	return attempts

}
