package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/cin-lawrence/gosandbox/pkg/api"
	"github.com/cin-lawrence/gosandbox/pkg/models"
	"github.com/cin-lawrence/gosandbox/pkg/services"
)

func init() {
	srv := services.NewUserService()

	_, err := srv.GetByEmail("surge@paragon.com")
	if err == nil {
		log.Infof("First user exists, skip creating...")
		return
	}

	userIn := models.UserInput{
		Name:     "Surge",
		Email:    "surge@paragon.com",
		Password: "paragon",
	}

	user, err := srv.Create(userIn)
	if err != nil {
		panic(err)
	}
	log.Infof("Seeded the first user %s (%s)", user.Name, user.Email)
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
