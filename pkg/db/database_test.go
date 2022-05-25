package db_test

import (
	"testing"

	"github.com/cin-lawrence/gosandbox/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestConnectDatabase(t *testing.T) {
	db.ConnectToDatabase()
	assert.NoError(t, nil)
}
