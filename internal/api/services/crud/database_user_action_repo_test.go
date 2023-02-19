package crud

import (
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/domain"
	"github.com/omalbum/omalbum-back/internal/api/utils/check"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func createDBUserActions() (*db.Database, func() error) {
	d := db.NewInMemoryDatabase()
	_ = d.Open()
	d.DB.LogMode(true)
	d.DB.CreateTable(&domain.User{})
	d.DB.CreateTable(&domain.UserAction{})

	return d, d.Close
}

func TestGetActionsByUserId(t *testing.T) {
	d, closeDb := createDBUserActions()
	defer check.Check(closeDb)

	userRepo := NewDatabaseUserRepo(d)
	d.DB.Create(&domain.User{
		UserName: "ivan",
		Email:    "ivansadofschi@gmail.com",
	})

	user := userRepo.GetByUserName("ivan")

	userActionRepo := NewDatabaseUserActionRepo(d)
	//timeNow:= string(time.Now())

	_ = userActionRepo.Create(&domain.UserAction{
		UserId:      user.ID,
		Date:        time.Now(),
		Description: "hizo una acción de prueba",
	})

	_ = userActionRepo.Create(&domain.UserAction{
		UserId:      user.ID,
		Date:        time.Now(),
		Description: "hizo otra accion de prueba",
	})

	userActions := userActionRepo.GetActionsByUserID(1)
	assert.Equal(t, 2, len(userActions))
	userAction := userActions[0]
	assert.Equal(t, user.ID, userAction.UserId)
	assert.Equal(t, uint(1), userAction.ID)
	//assert.Equal(t,timeNow, userAction.Date)
	assert.Equal(t, "hizo una acción de prueba", userAction.Description)

}

func TestGetById(t *testing.T) {
	d, closeDb := createDBUserActions()
	defer check.Check(closeDb)

	userRepo := NewDatabaseUserRepo(d)
	d.DB.Create(&domain.User{
		UserName: "ivan",
		Email:    "ivansadofschi@gmail.com",
	})
	user := userRepo.GetByUserName("ivan")

	userActionRepo := NewDatabaseUserActionRepo(d)

	ua := domain.UserAction{
		UserId:      user.ID,
		Date:        time.Now(),
		Description: "hizo una acción de prueba",
	}
	_ = userActionRepo.Create(&ua)

	userAction := userActionRepo.GetByID(ua.ID)
	assert.Equal(t, user.ID, userAction.UserId)
	assert.Equal(t, uint(1), userAction.ID)
	//assert.Equal(t,timeNow, userAction.Date)
	assert.Equal(t, "hizo una acción de prueba", userAction.Description)
	assert.Nil(t, userActionRepo.GetByID(12345))
}

func TestGetAll(t *testing.T) {
	d, closeDb := createDBUserActions()
	defer check.Check(closeDb)

	userRepo := NewDatabaseUserRepo(d)
	d.DB.Create(&domain.User{
		UserName: "ivan",
		Email:    "ivan@gmail.com",
	})
	d.DB.Create(&domain.User{
		UserName: "charly",
		Email:    "charly@gmail.com",
	})
	ivan := userRepo.GetByUserName("ivan")
	charly := userRepo.GetByUserName("charly")
	userActionRepo := NewDatabaseUserActionRepo(d)

	_ = userActionRepo.Create(&domain.UserAction{
		UserId:      ivan.ID,
		Date:        time.Now(),
		Description: "hizo una acción de prueba",
	})

	_ = userActionRepo.Create(&domain.UserAction{
		UserId:      charly.ID,
		Date:        time.Now(),
		Description: "hizo una acción de prueba",
	})

	userActions := userActionRepo.GetAll()
	assert.Equal(t, 2, len(userActions))

}
