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
	StatusCode  uint
	Method      string
	Resource    string
	Description string
}
