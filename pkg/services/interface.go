package services

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type EntityService struct {
        db *gorm.DB
}


func (srv *EntityService) newSession() *gorm.DB {
        ctx, _ := context.WithTimeout(context.Background(), time.Second)
        session := srv.db.WithContext(ctx)

        return session
}
