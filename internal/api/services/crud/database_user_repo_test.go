package crud

import (
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/domain"
	"github.com/omalbum/omalbum-back/internal/api/utils/check"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createDB() (*db.Database, func() error) {
	d := db.NewInMemoryDatabase()
	_ = d.Open()
	d.DB.LogMode(true)
	d.DB.CreateTable(&domain.User{})
	return d, d.Close
}

func TestGetUserByID(t *testing.T) {
	d, closeDb := createDB()
	defer check.Check(closeDb)

	userRepo := NewDatabaseUserRepo(d)
	d.DB.Create(&domain.User{
		UserName: "carlos",
	})

	user := userRepo.GetByID(1)
	assert.Equal(t, "carlos", user.UserName)

}

func TestCannotGetUserByID(t *testing.T) {
	d, closeDb := createDB()
	defer check.Check(closeDb)

	userRepo := NewDatabaseUserRepo(d)
	d.DB.Create(&domain.User{
		UserName: "carlos",
	})

	user := userRepo.GetByID(2)
	assert.Nil(t, user)
}

func TestGetUserByUserName(t *testing.T) {
	d, closeDb := createDB()
	defer check.Check(closeDb)

	userRepo := NewDatabaseUserRepo(d)
	d.DB.Create(&domain.User{
		UserName: "carlos",
	})

	user := userRepo.GetByUserName("carlos")
	assert.Equal(t, "carlos", user.UserName)
}

func TestGetUserByMail(t *testing.T) {
	d, closeDb := createDB()
	defer check.Check(closeDb)

	userRepo := NewDatabaseUserRepo(d)
	d.DB.Create(&domain.User{
		UserName: "carlos",
		Email:    "carlos@nada.com",
	})

	user := userRepo.GetByEmail("carlos@nada.com")
	assert.Equal(t, "carlos", user.UserName)
}

func TestCannotGetUserByUserName(t *testing.T) {
	d, closeDb := createDB()
	defer check.Check(closeDb)

	userRepo := NewDatabaseUserRepo(d)
	d.DB.Create(&domain.User{
		UserName: "carlos",
	})

	user := userRepo.GetByUserName("pepe")
	assert.Nil(t, user)
}

func TestCannotGetUserByEmail(t *testing.T) {
	d, closeDb := createDB()
	defer check.Check(closeDb)

	userRepo := NewDatabaseUserRepo(d)
	d.DB.Create(&domain.User{
		UserName: "carlos",
		Email:    "pepe@nada.com",
	})

	user := userRepo.GetByEmail("pepe")
	assert.Nil(t, user)
}

func TestCreateUserIsNotOk(t *testing.T) {
	d, closeDb := createDB()
	defer check.Check(closeDb)

	userRepo := NewDatabaseUserRepo(d)
	d.DB.Create(&domain.User{
		UserName: "carlos",
	})
	err := userRepo.Create(&domain.User{
		UserName: "carlos",
	})

	assert.NotNil(t, err)
}

func TestCreateUserIsOk(t *testing.T) {
	d, closeDb := createDB()
	defer check.Check(closeDb)

	userRepo := NewDatabaseUserRepo(d)
	err := userRepo.Create(&domain.User{
		UserName: "nuevouser",
	})

	assert.Nil(t, err)
}
