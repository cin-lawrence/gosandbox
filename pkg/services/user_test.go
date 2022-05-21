package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cin-lawrence/gosandbox/pkg/models"
	"github.com/cin-lawrence/gosandbox/pkg/services"
)

func TestUserService(t *testing.T) {
	srv := services.NewUserService()

	users, err := srv.List()
	assert.NoError(t, err)
	assert.Equal(t, len(users), 0)

	user := models.UserInput{
		Name:  "user-test",
		Email: "user-test@example.com",
	}

	userInDB, err := srv.Create(user)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, userInDB.Name)
	assert.Equal(t, user.Email, userInDB.Email)
	assert.NotEmpty(t, userInDB.CreatedAt)
	assert.NotEmpty(t, userInDB.UpdatedAt)
	assert.Equal(t, userInDB.CreatedAt, userInDB.UpdatedAt)

	usersInDB, err := srv.List()
	assert.NoError(t, err)
	assert.Equal(t, len(usersInDB), 1)

	userInDB = usersInDB[0]
	assert.Equal(t, user.Name, userInDB.Name)
	assert.Equal(t, user.Email, userInDB.Email)

	userInDB, err = srv.Get("1")
	assert.NoError(t, err)
	assert.Equal(t, user.Name, userInDB.Name)
	assert.Equal(t, user.Email, userInDB.Email)

	userUpdate := models.UserUpdate{
		Name: "user-test-updated",
	}
	userInDB, err = srv.Update("1", userUpdate)
	assert.NoError(t, err)
	assert.Equal(t, userUpdate.Name, userInDB.Name)
	assert.NotEqual(t, userInDB.CreatedAt, userInDB.UpdatedAt)

	err = srv.Delete("1")
	assert.NoError(t, err)

	userInDB, err = srv.Get("1")
	assert.Error(t, err)
}
