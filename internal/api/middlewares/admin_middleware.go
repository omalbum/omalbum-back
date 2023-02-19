package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/services/permissions"
	"net/http"
)

type AdminMiddleware interface {
	AdminCheck(c *gin.Context)
}

type adminMiddleware struct {
	database *db.Database
	manager  permissions.Manager
}

func NewAdminMiddleware(database *db.Database, manager permissions.Manager) AdminMiddleware {
	return &adminMiddleware{
		database: database,
		manager:  manager,
	}
}

func (a *adminMiddleware) AdminCheck(c *gin.Context) {
	if !a.manager.IsAdmin(c) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{})
		return
	}

	c.Next()
}
