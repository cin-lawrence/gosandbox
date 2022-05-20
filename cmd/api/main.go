package main

import (
	"os"

	"github.com/cin-lawrence/gosandbox/pkg/api"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// @title           Ninja REST API
// @version         1.0
// @description     This is a sample server celler server.

// @contact.name   Lawrence @ Cinnamon AI
// @contact.url    https://github.com/cin-lawrence
// @contact.email  lawrence@cinnamon.is

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl /api/v1/auth/login
func main() {
	server := api.NewAPIServer()
	server.Run("0.0.0.0", 8080)
}
