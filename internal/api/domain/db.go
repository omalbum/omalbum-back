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
	RegistrationDate time.Time
	LastActiveDate   time.Time
	Name             string
	LastName         string
	Cellphone        string // string por si alguien pone un numero que empieza con 00 o +
	Email            string `gorm:"unique;"`
	Notes            string
	IsAdmin          bool
}

type UserAction struct {
	gorm.Model
	UserId    uint
	Date      time.Time
	Action    string
	ExtraData string `gorm:"type:text"`
}
