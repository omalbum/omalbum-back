package permissions

import (
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/params"
)

type Manager interface {
	IsAdmin(context *gin.Context) bool
	IsAdminOrSameUser(context *gin.Context, userID uint) bool
}

type manager struct {
	database *db.Database
	userRepo domain.UserRepo
}

func NewManager(database *db.Database) Manager {
	return &manager{
		database: database,
		userRepo: crud.NewDatabaseUserRepo(database),
	}
}

func (p *manager) IsAdmin(context *gin.Context) bool {
	callerUserID := params.GetCallerID(context)
	return crud.NewDatabaseUserRepo(p.database).GetByID(callerUserID).IsAdmin
}

func (p *manager) IsAdminOrSameUser(context *gin.Context, userID uint) bool {
	callerUserID := params.GetCallerID(context)

	if callerUserID != userID && !p.IsAdmin(context) {
		return false
	}

	return true
}
