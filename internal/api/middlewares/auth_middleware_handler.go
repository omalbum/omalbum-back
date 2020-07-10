package middlewares

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/mailer"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/users"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/crypto"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/params"
	"net/http"
	"strings"
	"time"
)

type user struct {
	UserID   uint
	UserName string
}

type authMiddlewareHandler struct {
	database      *db.Database
	user          *user
	retrievedUser *domain.User
}

// Checks if the user is valid in the db
func (a *authMiddlewareHandler) authenticator(c *gin.Context) (interface{}, error) {
	var loginApp domain.LoginApp
	if err := c.ShouldBind(&loginApp); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	userRepo := crud.NewDatabaseUserRepo(a.database)

	err := domain.ValidateEmail(loginApp.UserName)
	if err == nil {
		a.retrievedUser = userRepo.GetByEmail(strings.ToLower(loginApp.UserName))
	} else {
		a.retrievedUser = userRepo.GetByUserName(strings.ToLower(loginApp.UserName))
	}

	// User do not exist
	if a.retrievedUser == nil {
		return nil, jwt.ErrFailedAuthentication
	}

	// Has valid password?
	if crypto.IsHashedPasswordEqualWithPlainPassword(a.retrievedUser.HashedPassword, loginApp.Password) {
		a.user = &user{
			UserID:   a.retrievedUser.ID,
			UserName: a.retrievedUser.UserName,
		}
		return a.user, nil
	}

	return nil, jwt.ErrFailedAuthentication
}

// Returns the payload for the token if the user is valid
func (a *authMiddlewareHandler) payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*user); ok {
		return jwt.MapClaims{
			params.IdentityKeyID: v.UserID,
			IdentityKeyUserName:  v.UserName,
		}
	}
	return jwt.MapClaims{}
}

// Once the token is received in a request and is valid, it extracts the data
func (a *authMiddlewareHandler) identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &user{
		UserID:   uint(claims[params.IdentityKeyID].(float64)),
		UserName: claims[IdentityKeyUserName].(string),
	}
}

// At this moment, the token is valid
// Saves user info in context
func (a *authMiddlewareHandler) authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*user); ok {
		c.Set(params.IdentityKeyID, v.UserID)
		c.Set(IdentityKeyUserName, v.UserName)
	}

	return true
}

// Message to show when the provided token is not valid
func (a *authMiddlewareHandler) unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func (a *authMiddlewareHandler) refreshResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, &domain.RefreshResponseApp{
		Token:      token,
		Expiration: expire.Format(time.RFC3339),
	})
}

func (a *authMiddlewareHandler) loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	user, _ := users.NewService(a.database, mailer.NewMock()).GetByUser(a.retrievedUser)

	c.JSON(http.StatusOK, &domain.LoginResponseApp{
		Token:      token,
		Expiration: expire.Format(time.RFC3339),
		User:       *user,
	})
}
