package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/check"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/crypto"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/testing_util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func createDBWithUser() (*db.Database, func() error) {
	d := db.NewInMemoryDatabase()
	_ = d.Open()
	d.DB.LogMode(true)
	d.DB.CreateTable(&domain.User{})

	userRepo := crud.NewDatabaseUserRepo(d)
	_ = userRepo.Create(&domain.User{
		UserName:       "admin",
		Email:          "admin@nada.com",
		HashedPassword: crypto.HashAndSalt("admin"),
	})

	return d, d.Close
}

func TestLoginOk(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUser)
	defer check.Check(closeDb)

	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"user_name\":\"admin\", \"password\":\"admin\"}"))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)

	NewAuthMiddleware(d).LoginHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginOkCaseInsensitive(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUser)
	defer check.Check(closeDb)

	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"user_name\":\"ADMIN\", \"password\":\"admin\"}"))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)

	NewAuthMiddleware(d).LoginHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginOkWithMail(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUser)
	defer check.Check(closeDb)

	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"user_name\":\"admin@nada.com\", \"password\":\"admin\"}"))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)

	NewAuthMiddleware(d).LoginHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginOkWithMailCaseInsensitive(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUser)
	defer check.Check(closeDb)

	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"user_name\":\"ADMIN@nada.com\", \"password\":\"admin\"}"))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)

	NewAuthMiddleware(d).LoginHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginNotOk(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUser)
	defer check.Check(closeDb)

	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"user_name\":\"admin\", \"password\":\"pepe\"}"))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)

	NewAuthMiddleware(d).LoginHandler(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginWithInvalidUser(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUser)
	defer check.Check(closeDb)

	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"user_name\":\"adminpepe\", \"password\":\"pepe\"}"))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)

	NewAuthMiddleware(d).LoginHandler(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMiddlewareWithOkToken(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUser)
	defer check.Check(closeDb)

	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"user_name\":\"admin\", \"password\":\"admin\"}"))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)

	NewAuthMiddleware(d).LoginHandler(c)

	var authInfo struct {
		Code   int       `json:"code"`
		Expire time.Time `json:"expire"`
		Token  string    `json:"token"`
	}

	_ = json.Unmarshal(w.Body.Bytes(), &authInfo)

	// Send token
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)

	c2.Request, _ = http.NewRequest("GET", "/", nil)
	c2.Request.Header.Add("Authorization", "Bearer "+authInfo.Token)

	NewAuthMiddleware(d).MiddlewareFunc()(c2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestMiddlewareWithNotOkToken(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUser)
	defer check.Check(closeDb)

	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Add("Authorization", "Bearer BADTOKEN")

	NewAuthMiddleware(d).MiddlewareFunc()(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMiddlewareRefreshTokenOK(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUser)
	defer check.Check(closeDb)

	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"user_name\":\"admin\", \"password\":\"admin\"}"))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)

	NewAuthMiddleware(d).LoginHandler(c)

	var authInfo struct {
		Code   int       `json:"code"`
		Expire time.Time `json:"expire"`
		Token  string    `json:"token"`
	}

	_ = json.Unmarshal(w.Body.Bytes(), &authInfo)

	// Send token
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)

	c2.Request, _ = http.NewRequest("GET", "/", nil)
	c2.Request.Header.Add("Authorization", "Bearer "+authInfo.Token)

	NewAuthMiddleware(d).RefreshHandler(c2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestMiddlewareRefreshTokenNotOk(t *testing.T) {
	d, closeDb, w, c := testing_util.SetupWithDBR(createDBWithUser)
	defer check.Check(closeDb)

	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Add("Authorization", "Bearer BADTOKEN")

	NewAuthMiddleware(d).RefreshHandler(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
