package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cin-lawrence/gosandbox/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestLiveness(t *testing.T) {
	api := api.NewAPIServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	api.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReadiness(t *testing.T) {
	api := api.NewAPIServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz/readiness", nil)
	api.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestRedirectDocs(t *testing.T) {
	api := api.NewAPIServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/docs/", nil)
	api.Router.ServeHTTP(w, req)

	assert.Equal(t, 301, w.Code)
}
