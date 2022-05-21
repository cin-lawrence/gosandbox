package v1_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/cin-lawrence/gosandbox/pkg/api/v1"
	"github.com/cin-lawrence/gosandbox/pkg/models"
	"github.com/cin-lawrence/gosandbox/pkg/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var userTokens = generateTokens()

func setupTestV1UserAPIRouter() *gin.Engine {
	router := gin.New()
	rg := router.Group("/")
	v1.NewV1UserGroup(rg)

	userForm := url.Values{
		"username": []string{"user-auth-login@example.com"},
		"password": []string{"user-test-password"},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(userForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	var tokens models.Tokens
	json.Unmarshal([]byte(w.Body.String()), &tokens)

	return router
}

func generateTokens() models.Tokens {
	router := setupTestV1AuthAPIRouter()

	srv := services.NewUserService()
	userIn := models.UserInput{
		Name:     "admin",
		Email:    "admin@example.com",
		Password: "user-test-password",
	}
	srv.Create(userIn)
	userForm := url.Values{
		"username": []string{"user-auth-login@example.com"},
		"password": []string{"user-test-password"},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(userForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	var tokens models.Tokens
	json.Unmarshal([]byte(w.Body.String()), &tokens)
	return tokens
}

func setAccessToken(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userTokens.AccessToken))
}

func TestV1UserAPIList(t *testing.T) {
	router := setupTestV1UserAPIRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/", nil)
	setAccessToken(req)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var users models.UserList
	err := json.Unmarshal([]byte(w.Body.String()), &users)
	assert.Nil(t, err)
	// assert.Equal(t, 0, len(users.Items))
}

func TestV1UserAPICreate(t *testing.T) {
	router := setupTestV1UserAPIRouter()
	userForm := url.Values{
		"name":     []string{"user-test"},
		"email":    []string{"user-test@example.com"},
		"password": []string{"user-test-password"},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/", strings.NewReader(userForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	setAccessToken(req)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var user models.User
	err := json.Unmarshal([]byte(w.Body.String()), &user)
	assert.Nil(t, err)
	assert.Equal(t, user.Name, "user-test")
	assert.Equal(t, user.Email, "user-test@example.com")
}

func TestV1UserAPIGet(t *testing.T) {
	router := setupTestV1UserAPIRouter()

	srv := services.NewUserService()
	userIn := models.UserInput{
		Name:     "user-test",
		Email:    "user-test-get@example.com",
		Password: "user-test-password",
	}
	user, err := srv.Create(userIn)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%v", user.ID), nil)
	setAccessToken(req)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal([]byte(w.Body.String()), &user)
	assert.Nil(t, err)
	assert.Equal(t, user.Name, "user-test")
	assert.Equal(t, user.Email, "user-test-get@example.com")
}

func TestV1UserAPIUpdate(t *testing.T) {
	router := setupTestV1UserAPIRouter()

	srv := services.NewUserService()
	userIn := models.UserInput{
		Name:     "user-test",
		Email:    "user-test-update@example.com",
		Password: "user-test-password",
	}
	user, err := srv.Create(userIn)
	assert.Nil(t, err)

	userUpd := models.UserUpdate{
		Name: "user-test-updated",
	}
	w := httptest.NewRecorder()
	jsonBody, _ := json.Marshal(userUpd)
	requestBody := bytes.NewBuffer(jsonBody)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%v", user.ID), requestBody)
	setAccessToken(req)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal([]byte(w.Body.String()), &user)
	assert.Nil(t, err)
	assert.Equal(t, user.Name, "user-test-updated")
	assert.Equal(t, user.Email, "user-test-update@example.com")
}

func TestV1UserAPIDelete(t *testing.T) {
	router := setupTestV1UserAPIRouter()

	srv := services.NewUserService()
	userIn := models.UserInput{
		Name:     "user-test",
		Email:    "user-test-delete@example.com",
		Password: "user-test-password",
	}
	user, err := srv.Create(userIn)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%v", user.ID), nil)
	setAccessToken(req)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}
