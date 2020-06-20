package middlewares

import (
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/config"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"log"
	"time"
)

const (
	IdentityKeyUserName = "user_name"
	realm               = "private zone"
)

type AuthMiddleware interface {
	LoginHandler(c *gin.Context)
	MiddlewareFunc() gin.HandlerFunc
	RefreshHandler(c *gin.Context)
}

type authMiddleware struct {
	ginJWTMiddleware *jwt.GinJWTMiddleware
}

func NewAuthMiddleware(database *db.Database) AuthMiddleware {
	authMiddlewareHandler := authMiddlewareHandler{
		database: database,
	}

	ginJWTMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           realm,
		Key:             []byte(config.AuthKey),
		Timeout:         24 * time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     IdentityKeyUserName,
		PayloadFunc:     authMiddlewareHandler.payloadFunc,
		IdentityHandler: authMiddlewareHandler.identityHandler,
		Authenticator:   authMiddlewareHandler.authenticator,
		Authorizator:    authMiddlewareHandler.authorizator,
		Unauthorized:    authMiddlewareHandler.unauthorized,
		LoginResponse:   authMiddlewareHandler.loginResponse,
		RefreshResponse: authMiddlewareHandler.refreshResponse,
		TokenLookup:     "header: Authorization",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return &authMiddleware{
		ginJWTMiddleware: ginJWTMiddleware,
	}
}

func (a *authMiddleware) LoginHandler(c *gin.Context) {
	a.ginJWTMiddleware.LoginHandler(c)
}

func (a *authMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return a.ginJWTMiddleware.MiddlewareFunc()
}

func (a *authMiddleware) RefreshHandler(c *gin.Context) {
	a.ginJWTMiddleware.RefreshHandler(c)
}
