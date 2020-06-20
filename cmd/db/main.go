package main

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/config"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/check"
	"log"
)

// This programs creates all the tables in the database based on the domain models
// Use carefully
func main() {
	database := db.NewServerDatabase(config.GetDatabase())

	err := database.Open()
	if err != nil {
		log.Print(err)
		return
	}
	defer check.Check(database.Close)

	log.Print("Tables are going to be created")
	database.DB.LogMode(true)

	crud.DropTables(database)
	crud.CreateTables(database)
	crud.CreateForeignKeys(database)
	crud.CreateAnonymousUser(database) // this user should be created first in order to have ID=1
	crud.CreateAdminUser(database)
	crud.CreateSampleData(database)
	log.Print("Tables created")
}
