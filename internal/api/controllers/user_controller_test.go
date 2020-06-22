package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/mailer"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/permissions"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/check"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/crypto"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/params"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/testing_util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func createDBWithUserAndCook() (*db.Database, func() error) {
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

	return database, database.Close
}

func TestCanGetUserOwnInfo(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUserAndCook)
	defer check.Check(closeDb)

	c.Request, _ = http.NewRequest("GET", "", nil)
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)
	c.Params = append(c.Params, gin.Param{
		Key:   "user_id",
		Value: "1",
	})
	c.Set(params.IdentityKeyID, uint(1))

	NewUserController(d, permissions.NewManager(d), mailer.NewMock()).GetUser(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCannotGetUserBecauseItIsNotFound(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUserAndCook)
	defer check.Check(closeDb)

	c.Request, _ = http.NewRequest("GET", "", nil)
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)
	c.Params = append(c.Params, gin.Param{
		Key:   "user_id",
		Value: "999",
	})
	c.Set(params.IdentityKeyID, uint(2))

	NewUserController(d, permissions.NewManager(d), mailer.NewMock()).GetUser(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCannotGetUserOthersInfoBecauseItIsNotAdmin(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUserAndCook)
	defer check.Check(closeDb)

	c.Request, _ = http.NewRequest("GET", "", nil)
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)
	c.Params = append(c.Params, gin.Param{
		Key:   "user_id",
		Value: "2",
	})
	c.Set(params.IdentityKeyID, uint(1))

	NewUserController(d, permissions.NewManager(d), mailer.NewMock()).GetUser(c)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCanGetOtherUserInfoBecauseItIsAdmin(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUserAndCook)
	defer check.Check(closeDb)

	c.Request, _ = http.NewRequest("GET", "", nil)
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)
	c.Params = append(c.Params, gin.Param{
		Key:   "user_id",
		Value: "1",
	})
	c.Set(params.IdentityKeyID, uint(2))

	NewUserController(d, permissions.NewManager(d), mailer.NewMock()).GetUser(c)
	assert.Equal(t, http.StatusOK, w.Code)
}
