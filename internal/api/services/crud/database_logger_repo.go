package crud

import (
	"encoding/json"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"time"
)

/*
	To log user's actions
	actionExtraData must be a string or something that can be Marshalled
*/
type Logger interface {
	LogUserAction(userID uint, actionName string, actionExtraData interface{}) error
	LogUserActionAtTime(userID uint, actionName string, actionExtraData interface{}, actionTime time.Time) error
	LogAnonymousAction(actionName string, actionExtraData interface{}) error
	LogAnonymousActionAtTime(actionName string, actionExtraData interface{}, actionTime time.Time) error
}

// Implementation

type logger struct {
	userActionRepo domain.UserActionRepo
}

func NewLogger(database *db.Database) *logger {
	return &logger{
		userActionRepo: NewDatabaseUserActionRepo(database),
	}
}

func (l *logger) LogUserActionAtTime(userID uint, actionName string, actionExtraData interface{}, actionTime time.Time) {
	s, ok := actionExtraData.(string)
	if !ok {
		s1, _ := json.Marshal(actionExtraData)
		s = string(s1)
	}
	_ = l.userActionRepo.Create(&domain.UserAction{
		UserId:    userID,
		Date:      actionTime,
		Action:    actionName,
		ExtraData: s,
	})
}

func (l *logger) LogAnonymousActionAtTime(actionName string, actionExtraData interface{}, actionTime time.Time) {
	l.LogUserActionAtTime(domain.AnonymousUser, actionName, actionExtraData, actionTime)
}

func (l *logger) LogUserAction(userID uint, actionName string, actionExtraData interface{}) {
	l.LogUserActionAtTime(userID, actionName, actionExtraData, time.Now())
}

func (l *logger) LogAnonymousAction(actionName string, actionExtraData interface{}) {
	l.LogAnonymousActionAtTime(actionName, actionExtraData, time.Now())
}
