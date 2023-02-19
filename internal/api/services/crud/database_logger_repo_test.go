package crud

import (
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/domain"
	"github.com/omalbum/omalbum-back/internal/api/utils/check"
	"github.com/stretchr/testify/assert"
	"net/http"
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

	newLogger.LogAnonymousAction("tried to register", http.StatusBadRequest, "GET", "/register")

	newLogger.LogUserAction(ivan.ID, "hizo una acción de prueba", http.StatusForbidden, "POST", "/answer")

	userActions := userActionRepo.GetActionsByUserID(ivan.ID)
	assert.Equal(t, 1, len(userActions))
	userAction := userActions[0]
	assert.Equal(t, "hizo una acción de prueba", userAction.Description)
	assert.Equal(t, "/answer", userAction.Resource)

	userActions = userActionRepo.GetActionsByUserID(domain.AnonymousUser)
	userAction = userActions[0]
	assert.Equal(t, "tried to register", userAction.Description)
	assert.Equal(t, "/register", userAction.Resource)
	assert.Equal(t, http.StatusBadRequest, userAction.StatusCode)
}
