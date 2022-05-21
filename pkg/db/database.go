package db

import (
	"errors"
	"strings"

	"github.com/cin-lawrence/gosandbox/pkg/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var DB *gorm.DB = ConnectToDatabase()

func ConnectToDatabase() *gorm.DB {
	var dbConn *gorm.DB
	var err error

	switch {
	case strings.HasPrefix(config.Config.DatabaseURI, "postgresql://"):
		dbConn, err = ConnectToPostgres(config.Config.DatabaseURI)
	case strings.HasPrefix(config.Config.DatabaseURI, "sqlite://"):
		dbConn, err = ConnectToSqlite(config.Config.DatabaseURI)
	default:
		err = errors.New("Unsupported DB URI")
	}

	if err != nil {
		log.Fatal("Failed to connect database")
		panic(err)
	}

	return dbConn
}
