package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
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
	cache    domain.TeleOMACache
}

func NewUserProblemMiddleware(database *db.Database, manager permissions.Manager, cache domain.TeleOMACache) UserProblemMiddleware {
	return &userProblemMiddleware{
		database: database,
		manager:  manager,
		cache:    cache,
	}
}

func (a *userProblemMiddleware) ViewAuthCheck(c *gin.Context) {
	problemId := params.GetProblemID(c)
	key := domain.ProblemViewableCacheKey(problemId)
	var res = a.cache.Get(key)
	var problemIsViewable bool
	if res != nil {
		problemIsViewable = res.(bool)
	} else {
		problem := crud.NewDatabaseProblemRepo(a.database).GetById(problemId)
		if problem == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
			return
		}
		problemIsViewable = problem.IsViewable()
		a.cache.SetWithTTL(key, problemIsViewable, domain.DefaultTimeToLive)
	}
	if !problemIsViewable {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{})
		return
	}
	c.Next()
}
