package db

import (
        "gorm.io/driver/postgres"
        "gorm.io/gorm"
)

func ConnectToPostgres(uri string) (*gorm.DB, error) {
        return gorm.Open(postgres.Open(uri), &gorm.Config{})
}
