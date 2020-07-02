package permissions

import (
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/check"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/crypto"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/params"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/testing_util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createDBWithAdmin() (*db.Database, func() error) {
	database := db.NewInMemoryDatabase()
	_ = database.Open()
	database.DB.LogMode(true)

	crud.CreateTables(database)

	userRepo := crud.NewDatabaseUserRepo(database)
	user := domain.User{
		UserName:       "cook",
		HashedPassword: crypto.HashAndSalt("pastafrola"),
		Email:          "el_cook@gmail.com",
	}
	_ = userRepo.Create(&user)

	userAdmin := domain.User{
		UserName:       "gorra",
		HashedPassword: crypto.HashAndSalt("gorra"),
		Email:          "gorra15@gmail.com",
		IsAdmin:        true,
	}
	_ = userRepo.Create(&userAdmin)

	userOther := domain.User{
		UserName:       "other",
		HashedPassword: crypto.HashAndSalt("gorra"),
		Email:          "other@gmail.com",
	}
	_ = userRepo.Create(&userOther)

	return database, database.Close
}

func TestIsAdmin(t *testing.T) {
	d, closeDb, c := testing_util.SetupWithDB(createDBWithAdmin)
	defer check.Check(closeDb)

	c.Set(params.IdentityKeyID, uint(2))

	assert.True(t, NewManager(d).IsAdmin(c))
}

func TestIsAdminOrSameUser(t *testing.T) {
	d, closeDb, c := testing_util.SetupWithDB(createDBWithAdmin)
	defer check.Check(closeDb)

	c.Params = append(c.Params, gin.Param{
		Key:   "user_id",
		Value: "1",
	})
	c.Set(params.IdentityKeyID, uint(1))

	assert.True(t, NewManager(d).IsAdminOrSameUser(c, 1))
}

func TestIsNotAdminOrSameUser(t *testing.T) {
	d, closeDb, c := testing_util.SetupWithDB(createDBWithAdmin)
	defer check.Check(closeDb)

	c.Params = append(c.Params, gin.Param{
		Key:   "user_id",
		Value: "3",
	})
	c.Set(params.IdentityKeyID, uint(3))

	assert.False(t, NewManager(d).IsAdminOrSameUser(c, 1))
}
