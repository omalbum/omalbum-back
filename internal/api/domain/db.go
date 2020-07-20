package domain

import (
	"github.com/jinzhu/gorm"
	"time"
)

const (
	AnonymousUser uint = 1
)

// Si agregamos algun campo a este struct hay que tener mucho cuidado
// de no pisarlo con null en el endpoint de Update user!
type User struct {
	gorm.Model
	// campos de registro de usuario
	UserName       string `gorm:"unique;size:20"`
	HashedPassword string
	Email          string `gorm:"unique;"`
	Name           string
	LastName       string
	BirthDate      time.Time
	Gender         string
	Country        string
	Province       string
	Department     string
	Location       string
	School         string
	IsStudent      bool
	IsTeacher      bool
	SchoolYear     uint
	// campos que se llenan de otra forma
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
	DateContestStart time.Time
	DateContestEnd   time.Time
	Statement        string `gorm:"type:longtext;"`
	Answer           int
	Annotations      string `gorm:"type:longtext;"`
	IsDraft          bool
	Hint             string `gorm:"type:longtext;"`
	OfficialSolution string `gorm:"type:longtext;"`
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
