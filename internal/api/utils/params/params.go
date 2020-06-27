package params

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

const IdentityKeyID = "ID"

// Generic function to retrieve a param
func GetUintValueFromParam(context *gin.Context, key string) uint {
	value, _ := strconv.ParseInt(context.Param(key), 10, 32)
	return uint(value)
}

func GetStringValueFromParam(context *gin.Context, key string) string {
	value := context.Param(key)
	return value
}

func GetBoolValueFromParam(context *gin.Context, key string) bool {
	value, _ := strconv.ParseBool(context.Param(key))
	return value
}

// Returns the user id from the params
func GetUserID(context *gin.Context) uint {
	return GetUintValueFromParam(context, "user_id")
}

func GetProblemID(context *gin.Context) uint {
	return GetUintValueFromParam(context, "problem_id")
}

// Returns the caller id, extracted from the token
func GetCallerID(context *gin.Context) uint {
	IdentityKeyID, _ := context.Get(IdentityKeyID)
	callerUserID, _ := IdentityKeyID.(uint)

	return callerUserID
}
