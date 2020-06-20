package params

import (
	"github.com/gin-gonic/gin"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/testing_util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserID(t *testing.T) {
	c := testing_util.Setup()

	c.Params = append(c.Params, gin.Param{
		Key:   "user_id",
		Value: "2",
	})

	assert.Equal(t, uint(2), GetUserID(c))
}

func TestGetCallerID(t *testing.T) {
	c := testing_util.Setup()

	c.Set(IdentityKeyID, uint(100))

	assert.Equal(t, uint(100), GetCallerID(c))
}
