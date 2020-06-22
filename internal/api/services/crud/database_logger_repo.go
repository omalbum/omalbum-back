package crud

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"time"
)

/*
	To log user's actions
*/
type Logger interface {
	LogUserAction(userID uint, description string, statusCode uint, method string, resource string) error
	LogUserActionAtTime(userID uint, actionTime time.Time, description string, statusCode uint, method string, resource string) error
	LogAnonymousAction(description string, statusCode uint, method string, resource string) error
	LogAnonymousActionAtTime(actionTime time.Time, description string, statusCode uint, method string, resource string) error
}

// Implementation

type logger struct {
	userActionRepo domain.UserActionRepo
}

func NewLogger(database *db.Database) Logger {
	return &logger{
		userActionRepo: NewDatabaseUserActionRepo(database),
	}
}

func (l *logger) LogUserActionAtTime(userID uint, actionTime time.Time, description string, statusCode uint, method string, resource string) error {
	return l.userActionRepo.Create(&domain.UserAction{
		UserId:      userID,
		Date:        actionTime,
		StatusCode:  statusCode,
		Method:      method,
		Resource:    resource,
		Description: description,
	})
}

func (l *logger) LogAnonymousActionAtTime(actionTime time.Time, description string, statusCode uint, method string, resource string) error {
	return l.LogUserActionAtTime(domain.AnonymousUser, actionTime, description, statusCode, method, resource)
}

func (l *logger) LogUserAction(userID uint, description string, statusCode uint, method string, resource string) error {
	return l.LogUserActionAtTime(userID, time.Now(), description, statusCode, method, resource)
}

func (l *logger) LogAnonymousAction(description string, statusCode uint, method string, resource string) error {
	return l.LogAnonymousActionAtTime(time.Now(), description, statusCode, method, resource)
}
