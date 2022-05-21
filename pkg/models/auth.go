package models

import (
	uuid "github.com/satori/go.uuid"
)

type RefreshTokenInput struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserLogin struct {
	Username string `form:"username" json:"username" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=50"`
}

type TokenMeta struct {
	AccessToken    string
	RefreshToken   string
	AccessUUID     uuid.UUID
	RefreshUUID    uuid.UUID
	AccessExpires  int64
	RefreshExpires int64
}

type AccessMeta struct {
	AccessUUID uuid.UUID
	UserID     int64
}
