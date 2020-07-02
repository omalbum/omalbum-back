package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/mailer"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/check"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/crypto"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/testing_util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func createDBWithUser() (*db.Database, func() error) {
	database := db.NewInMemoryDatabase()
	_ = database.Open()
	database.DB.LogMode(true)

	crud.CreateTables(database)

	userRepo := crud.NewDatabaseUserRepo(database)
	_ = userRepo.Create(&domain.User{
		UserName:       "admin",
		HashedPassword: crypto.HashAndSalt("admin"),
	})

	return database, database.Close
}

func TestCanCreateUser(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUser)
	defer check.Check(closeDb)

	request := `
{
    "user_name": "cdesseno2",
    "password": "password123",
    "name": "Carlos",
    "last_name": "Desseno",
	"birthdate": "1992-03-14T00:00:00-00:00",
    "email": "cdesseno@gmail.com"
}
`
	c.Request, _ = http.NewRequest("POST", "", bytes.NewBufferString(request))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)

	NewRegisterController(d, mailer.NewMock()).Register(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCannotCreateUserBecauseThereIsAnotherOneWithTheSameUserName(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUser)
	defer check.Check(closeDb)

	request := `
{
    "user_name": "admin",
    "password": "password123",
    "name": "Carlos",
    "last_name": "Desseno",
	"birthdate": "1992-03-14T00:00:00-00:00",
    "email": "cdesseno@gmail.com"
}
`
	c.Request, _ = http.NewRequest("POST", "", bytes.NewBufferString(request))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)

	NewRegisterController(d, mailer.NewMock()).Register(c)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestCannotCreateUserBecauseValidationFails(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUser)
	defer check.Check(closeDb)

	request := `
{
    "user_name": "admin_otra_pass",
    "password": "pass",
    "name": "Carlos",
    "last_name": "Desseno",
    "cellphone": "1233333",
    "email": "cdesseno@gmail.com",
}
`
	c.Request, _ = http.NewRequest("POST", "", bytes.NewBufferString(request))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)

	NewRegisterController(d, mailer.NewMock()).Register(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
