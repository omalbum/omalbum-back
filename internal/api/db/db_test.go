package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenDatabaseIsOk(t *testing.T) {
	db := NewInMemoryDatabase()
	err := db.Open()

	assert.Nil(t, err)
}

func TestCloseDatabaseIsOk(t *testing.T) {
	db := NewInMemoryDatabase()
	_ = db.Open()
	err := db.Close()

	assert.Nil(t, err)
}

func TestOpenDatabaseFails(t *testing.T) {
	db := NewFakeServerDatabase()
	err := db.Open()

	assert.NotNil(t, err)
}

// TODO
func TestCloseDatabaseFails(t *testing.T) {

}
