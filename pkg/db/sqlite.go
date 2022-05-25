package db

import (
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectToSqlite(uri string) (*gorm.DB, error) {
	uri = strings.Replace(uri, "sqlite://", "", 1)
	return gorm.Open(sqlite.Open(uri), &gorm.Config{})
}
