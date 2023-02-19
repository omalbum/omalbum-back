package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/domain"
	"github.com/omalbum/omalbum-back/internal/api/services/crud"
	"github.com/omalbum/omalbum-back/internal/api/services/permissions"
	"github.com/omalbum/omalbum-back/internal/api/utils/check"
	"github.com/omalbum/omalbum-back/internal/api/utils/params"
	"github.com/omalbum/omalbum-back/internal/api/utils/testing_util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func createDBForGroupsTesting() (*db.Database, func() error) {
	database := db.NewInMemoryDatabase()
	_ = database.Open()
	database.DB.LogMode(true)

	crud.CreateTables(database)

	userRepo := crud.NewDatabaseUserRepo(database)
	_ = userRepo.Create(&domain.User{
		UserName: "admin",
		Email:    "a@b.com",
		IsAdmin:  true,
	})
	_ = userRepo.Create(&domain.User{
		UserName: "cook",
		Email:    "c@b.com",
	})

	return database, database.Close
}

func TestItIsNotAdmin(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBForGroupsTesting)
	defer check.Check(closeDb)

	request := `{}`
	c.Request, _ = http.NewRequest("POST", "", bytes.NewBufferString(request))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)
	c.Params = append(c.Params, gin.Param{
		Key:   "user_id",
		Value: "2",
	})
	c.Set(params.IdentityKeyID, uint(2))

	NewAdminMiddleware(d, permissions.NewManager(d)).AdminCheck(c)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestItIsAdmin(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBForGroupsTesting)
	defer check.Check(closeDb)

	request := `{}`
	c.Request, _ = http.NewRequest("POST", "", bytes.NewBufferString(request))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)
	c.Params = append(c.Params, gin.Param{
		Key:   "user_id",
		Value: "1",
	})
	c.Set(params.IdentityKeyID, uint(1))

	NewAdminMiddleware(d, permissions.NewManager(d)).AdminCheck(c)
	assert.Equal(t, http.StatusOK, w.Code)
}
