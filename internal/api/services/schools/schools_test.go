package schools

import (
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/domain"
	"github.com/omalbum/omalbum-back/internal/api/services/crud"
	"github.com/omalbum/omalbum-back/internal/api/utils/check"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createDBWithSchool() (*db.Database, func() error) {
	database := db.NewInMemoryDatabase()
	_ = database.Open()
	database.DB.LogMode(true)

	crud.CreateTables(database)

	schoolRepo := crud.NewDatabaseSchoolRepo(database)
	school := domain.School{
		Name:       "Escuela N° 1",
		Province:   "Buenos Aires",
		Department: "Depto",
	}
	_ = schoolRepo.Create(&school)

	return database, database.Close
}

func TestCanGetBySchoolName(t *testing.T) {
	d, closeDb := createDBWithSchool()
	defer check.Check(closeDb)

	schoolApp := NewService(d).GetSchools("cuela", "Buenos Aires", "Depto")
	assert.NotNil(t, schoolApp)
}

func TestCannotGetSchoolBecauseItIsNotFound(t *testing.T) {
	d, closeDb := createDBWithSchool()
	defer check.Check(closeDb)

	schools := NewService(d).GetSchools("escuela", "asdasd", "ooooo")
	assert.Equal(t, 0, len(schools))
}
