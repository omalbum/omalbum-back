package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashAndSaltIsOk(t *testing.T) {
	hashedPassword := HashAndSalt("carlos")

	res := IsHashedPasswordEqualWithPlainPassword(hashedPassword, "carlos")
	assert.True(t, res)
}

func TestHashAndSaltIsNotOk(t *testing.T) {
	hashedPassword := HashAndSalt("carlos")

	res := IsHashedPasswordEqualWithPlainPassword(hashedPassword, "ivan")
	assert.False(t, res)
}
