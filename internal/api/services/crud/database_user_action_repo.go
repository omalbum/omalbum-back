package crud

import (
	"github.com/jinzhu/gorm"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
)

type databaseUserActionRepo struct {
	database *db.Database
}

func NewDatabaseUserActionRepo(database *db.Database) domain.UserActionRepo {
	return &databaseUserActionRepo{
		database: database,
	}
}

func (dr *databaseUserActionRepo) GetByID(ID uint) *domain.UserAction {
	var userAction domain.UserAction
	if dr.database.DB.Where(&domain.UserAction{Model: gorm.Model{ID: ID}}).First(&userAction).RecordNotFound() {
		return nil
	}
	return &userAction
}

// Returns all actions for a user
func (dr *databaseUserActionRepo) GetActionsByUserID(userId uint) []domain.UserAction {
	var userActions []domain.UserAction
	dr.database.DB.Where(&domain.UserAction{UserId: userId}).Find(&userActions)
	return userActions
}

// returns all actions
func (dr *databaseUserActionRepo) GetAll() []domain.UserAction {
	var userActions []domain.UserAction
	dr.database.DB.Find(&userActions)

	return userActions
}

// Creates a userAction
func (dr *databaseUserActionRepo) Create(userAction *domain.UserAction) error {
	return dr.database.DB.Create(userAction).Error
}
