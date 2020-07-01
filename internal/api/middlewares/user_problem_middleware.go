package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/permissions"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/params"
	"net/http"
)

type UserProblemMiddleware interface {
	ViewAuthCheck(c *gin.Context)
}

type userProblemMiddleware struct {
	database *db.Database
	manager  permissions.Manager
}

func NewUserProblemMiddleware(database *db.Database, manager permissions.Manager) UserProblemMiddleware {
	return &userProblemMiddleware{
		database: database,
		manager:  manager,
	}
}

func (a *userProblemMiddleware) ViewAuthCheck(c *gin.Context) {

	problemId := params.GetProblemID(c)
	problem := crud.NewDatabaseProblemRepo(a.database).GetById(problemId)

	if !problem.IsViewable() {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{})
		return
	}
	c.Next()
}
