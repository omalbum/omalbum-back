package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/miguelsotocarlos/teleoma/internal/api/config"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"log"
)

type Database struct {
	DB      *gorm.DB
	dialect string
	args    string
}

func NewInMemoryDatabase() *Database {
	return &Database{
		dialect: "sqlite3",
		args:    ":memory:?cache=shared",
	}
}

func NewFakeServerDatabase() *Database {
	return &Database{
		dialect: "mysql",
		args:    "fake:fake@(localhost)/fake?charset=utf8&parseTime=True&loc=Local",
	}
}

func NewServerDatabase(c config.Database) *Database {
	return &Database{
		dialect: "mysql",
		args:    c.User + ":" + c.Password + "@(" + c.Host + ")/" + c.Database + "?charset=utf8&parseTime=True&loc=Local",
	}
}

func (d *Database) Open() error {
	log.Print("Opening database...")
	db, err := gorm.Open(d.dialect, d.args)

	if err != nil {
		return messages.New("database_open_error", "something happened opening the database: "+err.Error())
	}

	// Auto preload associated entities: magic!!
	db = db.Set("gorm:auto_preload", true)

	d.DB = db
	return nil
}

func (d *Database) Close() error {
	err := d.DB.Close()

	if err != nil {
		return messages.New("database_close_error", "something happened closing the database: "+err.Error())
	}

	log.Print("Closing database...")
	return nil
}
