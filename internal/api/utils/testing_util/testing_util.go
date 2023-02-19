package testing_util

import (
	"github.com/gin-gonic/gin"
	"github.com/omalbum/omalbum-back/internal/api/db"
	"net/http/httptest"
)

func SetupWithDBR(dbFunc func() (*db.Database, func() error)) (*db.Database, func() error, *httptest.ResponseRecorder, *gin.Context) {
	d, closeDb := dbFunc()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	return d, closeDb, w, c
}

func SetupWithDB(dbFunc func() (*db.Database, func() error)) (*db.Database, func() error, *gin.Context) {
	d, closeDb := dbFunc()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	return d, closeDb, c
}

func Setup() *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	return c
}
