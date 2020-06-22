package users

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/mailer"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/check"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createDBWithUserAndCook() (*db.Database, func() error) {
	database := db.NewInMemoryDatabase()
	_ = database.Open()
	database.DB.LogMode(true)

	crud.CreateTables(database)

	userRepo := crud.NewDatabaseUserRepo(database)
	user := domain.User{
		UserName:       "all",
		HashedPassword: crypto.HashAndSalt("pastafrola"),
		Email:          "el_cook@gmail.com",
	}
	_ = userRepo.Create(&user)

	return database, database.Close
}

func TestCanGetByUserID(t *testing.T) {
	d, closeDb := createDBWithUserAndCook()
	defer check.Check(closeDb)

	userApp, _ := NewService(d, mailer.NewMock()).GetByUserID(1)
	assert.NotNil(t, userApp)
}

func TestCannotGetUserBecauseItIsNotFound(t *testing.T) {
	d, closeDb := createDBWithUserAndCook()
	defer check.Check(closeDb)

	_, err := NewService(d, mailer.NewMock()).GetByUserID(99)
	assert.Equal(t, "user_not_found", messages.GetCode(err))
}

func TestCanGetByUser(t *testing.T) {
	d, closeDb := createDBWithUserAndCook()
	defer check.Check(closeDb)

	userApp, _ := NewService(d, mailer.NewMock()).GetByUser(crud.NewDatabaseUserRepo(d).GetByUserName("all"))
	assert.NotNil(t, userApp)
}

func TestServiceChangePassword(t *testing.T) {
	// testeamos que cambiar la password funcione
	// testeamos que las passwords de otros usuarios no sean pisadas
	database := db.NewInMemoryDatabase()
	_ = database.Open()
	database.DB.LogMode(true)

	crud.CreateTables(database)

	userRepo := crud.NewDatabaseUserRepo(database)
	aCook := domain.User{
		UserName:       "elcook",
		HashedPassword: crypto.HashAndSalt("pastafrola"),
		Email:          "el_cook@gmail.com",
	}
	anotherCook := domain.User{
		UserName:       "martiniano",
		HashedPassword: crypto.HashAndSalt("casancrem"),
		Email:          "mmolina@gmail.com",
	}
	_ = userRepo.Create(&aCook)
	_ = userRepo.Create(&anotherCook)
	NewService(database, mailer.NewMock()).ChangePassword(aCook.ID, "chocotorta")
	aCook = *userRepo.GetByID(aCook.ID)
	assert.True(t, crypto.IsHashedPasswordEqualWithPlainPassword(aCook.HashedPassword, "chocotorta"))
	anotherCook = *userRepo.GetByID(anotherCook.ID)
	assert.True(t, crypto.IsHashedPasswordEqualWithPlainPassword(anotherCook.HashedPassword, "casancrem"))
	assert.False(t, crypto.IsHashedPasswordEqualWithPlainPassword(anotherCook.HashedPassword, "chocotorta"))
}

func TestServiceUpdateUser(t *testing.T) {
	database := db.NewInMemoryDatabase()
	_ = database.Open()
	database.DB.LogMode(true)

	crud.CreateTables(database)

	userRepo := crud.NewDatabaseUserRepo(database)
	aUser := domain.User{
		UserName:  "juancito",
		Name:      "juan",
		LastName:  "cito",
		Cellphone: "1234579",
		Email:     "juancito@gmail.com",
	}
	anotherUser := domain.User{
		UserName:  "martiniano",
		Name:      "Martiniano",
		LastName:  "Molina",
		Cellphone: "1234579",
		Email:     "mmolina@gmail.com",
	}
	_ = userRepo.Create(&aUser)
	_ = userRepo.Create(&anotherUser)
	s := NewService(database, mailer.NewMock())

	updatedProfile := domain.RegistrationApp{
		Name:     "Juan Martin",
		LastName: "Delpo",
	}
	s.UpdateUser(aUser.ID, &updatedProfile)
	aUser = *userRepo.GetByID(aUser.ID)
	assert.Equal(t, "Delpo", aUser.LastName)
	assert.Equal(t, "Juan Martin", aUser.Name)
	anotherUser = *userRepo.GetByID(anotherUser.ID)
	assert.Equal(t, "Molina", anotherUser.LastName)
	assert.Equal(t, "Martiniano", anotherUser.Name)

	updatedProfile.Email = "mmolina@gmail.com"
	_, err := s.UpdateUser(aUser.ID, &updatedProfile)
	assert.NotNil(t, err)

}

func TestServiceResetPassword(t *testing.T) {

	database := db.NewInMemoryDatabase()
	_ = database.Open()
	database.DB.LogMode(true)

	crud.CreateTables(database)

	userRepo := crud.NewDatabaseUserRepo(database)
	aUser := domain.User{
		UserName:  "juancito",
		Name:      "juan",
		LastName:  "cito",
		Cellphone: "1234579",
		Email:     "juancito@gmail.com",
	}
	_ = userRepo.Create(&aUser)
	s := NewService(database, mailer.NewMock())

	err := s.ResetPassword("juancito@gmail.com")
	assert.Nil(t, err)
	err = s.ResetPassword("juancito1@gmail.com")
	assert.NotNil(t, err)
}
