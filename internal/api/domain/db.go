package domain

import (
	"github.com/jinzhu/gorm"
	"time"
)

const (
	AnonymousUser uint = 1
)

type User struct {
	gorm.Model
	UserName         string `gorm:"unique;size:20"`
	HashedPassword   string
	Email            string `gorm:"unique;"`
	Name             string
	LastName         string
	BirthDate        time.Time
	Gender           string
	Country          string
	Province         string
	Department       string
	Location         string
	School           string
	IsStudent        bool
	SchoolYear       uint
	RegistrationDate time.Time
	LastActiveDate   time.Time
	IsAdmin          bool
}

type UserAction struct {
	gorm.Model
	UserId      uint
	Date        time.Time
	StatusCode  int
	Method      string
	Resource    string
	Description string
}

type Problem struct {
	gorm.Model
	PoserId          uint
	OmaforosPostId   uint
	DateUploaded     time.Time
	DateContestStart time.Time
	DateContestEnd   time.Time
	Statement        string
	Answer           int
	Annotations      string
	IsDraft          bool
	Hint             string
	OfficialSolution string
	Series           string
	NumberInSeries   uint
}

type UserProblemAttempt struct {
	gorm.Model
	UserId     uint
	ProblemId  uint
	Date       time.Time
	UserAnswer int
}

type ProblemTag struct {
	gorm.Model
	ProblemId uint
	Tag       string
}

type ExpandedUserProblemAttempt struct {
	Answer           int
	DateContestEnd   time.Time
	DateContestStart time.Time
	UserId           uint
	AttemptId        uint
	AttemptDate      time.Time
	ProblemId        uint
	UserAnswer       int
	IsCorrect        bool
	DuringContest    bool
}

type School struct {
	gorm.Model
	Name       string
	Province   string
	Department string
	Location   string
}
