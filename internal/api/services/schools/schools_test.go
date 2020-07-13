package schools

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/mailer"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/check"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createDBWithSchoolAndCook() (*db.Database, func() error) {
	database := db.NewInMemoryDatabase()
	_ = database.Open()
	database.DB.LogMode(true)

	crud.CreateTables(database)

	schoolRepo := crud.NewDatabaseSchoolRepo(database)
	school := domain.School{
		Name:       "Escuela N° 1",
		Province: "Buenos Aires",
		Department: "Depto",
	}
	_ = schoolRepo.Create(&school)

	return database, database.Close
}

func TestCanGetBySchoolName(t *testing.T) {
	d, closeDb := createDBWithSchoolAndCook()
	defer check.Check(closeDb)

	schoolApp, _ := NewService(d, mailer.NewMock()).GetSchools("cuela", "Buenos Aires", "Depto")
	assert.NotNil(t, schoolApp)
}

func TestCannotGetSchoolBecauseItIsNotFound(t *testing.T) {
	d, closeDb := createDBWithSchoolAndCook()
	defer check.Check(closeDb)

	_, err := NewService(d, mailer.NewMock()).GetSchools("escuela", "asdasd", "ooooo")
	assert.Equal(t, "school_not_found", messages.GetCode(err))
}