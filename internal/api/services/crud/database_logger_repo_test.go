package crud

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/check"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createDBLogger() (*db.Database, func() error) {
	database := db.NewInMemoryDatabase()
	_ = database.Open()
	database.DB.LogMode(true)
	database.DB.CreateTable(&domain.User{})
	database.DB.CreateTable(&domain.UserAction{})
	CreateAnonymousUser(database)
	return database, database.Close
}

func TestLogger(t *testing.T) {
	database, closeDb := createDBLogger()
	defer check.Check(closeDb)

	userRepo := NewDatabaseUserRepo(database)
	userActionRepo := NewDatabaseUserActionRepo(database)

	database.DB.Create(&domain.User{
		UserName: "ivan",
		Email:    "ivansadofschi@gmail.com",
	})

	ivan := userRepo.GetByUserName("ivan")

	newLogger := NewLogger(database)

	newLogger.LogAnonymousAction("tried to register", "{blabla1}")

	newLogger.LogUserAction(ivan.ID, "hizo una acción de prueba", "{blabla}")

	userActions := userActionRepo.GetActionsByUserID(ivan.ID)
	assert.Equal(t, 1, len(userActions))
	userAction := userActions[0]
	assert.Equal(t, "hizo una acción de prueba", userAction.Action)
	assert.Equal(t, "{blabla}", userAction.ExtraData)

	userActions = userActionRepo.GetActionsByUserID(domain.AnonymousUser)
	userAction = userActions[0]
	assert.Equal(t, "tried to register", userAction.Action)
	assert.Equal(t, "{blabla1}", userAction.ExtraData)
}
