package v1_test

import (
	"bytes"
	"encoding/json"
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

var accessToken string
var refreshToken string

func setupTestV1AuthAPIRouter() *gin.Engine {
	router := gin.New()
	rg := router.Group("/")
	v1.NewV1AuthGroup(rg)
	return router
}

func TestLogin(t *testing.T) {
	router := setupTestV1AuthAPIRouter()

	srv := services.NewUserService()
	userIn := models.UserInput{
		Name:     "user-test",
		Email:    "user-auth-login@example.com",
		Password: "user-test-password",
	}
	_, err := srv.Create(userIn)
	assert.Nil(t, err)
	userForm := url.Values{
		"username": []string{"user-auth-login@example.com"},
		"password": []string{"user-test-password"},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(userForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var tokens models.Tokens
	err = json.Unmarshal([]byte(w.Body.String()), &tokens)
	assert.Nil(t, err)

	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	accessToken = tokens.AccessToken
	refreshToken = tokens.RefreshToken
}

func TestRefreshToken(t *testing.T) {
	router := setupTestV1AuthAPIRouter()

	refreshPayload := models.RefreshTokenInput{
		RefreshToken: refreshToken,
	}
	data, _ := json.Marshal(refreshPayload)
	req, err := http.NewRequest("POST", "/auth/refresh", bytes.NewBufferString(string(data)))
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	t.Logf("%v", w)
	assert.Equal(t, http.StatusOK, w.Code)
}
