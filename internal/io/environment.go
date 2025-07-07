package io

import (
	"github.com/caarlos0/env/v11"
	"github.com/motojouya/geezer_auth/internal/core/db"
	"github.com/motojouya/geezer_auth/pkg/core/jwt"
)

type Environment interface {
	GetJwtHandling() (jwt.JwtHandling, error)
	GetDBAccess() (db.DBAccess, error)
}

type environment struct{}

func CreateEnvironment() Environment {
	return &environment{}
}

func (e environment) GetJwtHandling() (jwt.JwtHandling, error) {
	return env.ParseAs[jwt.JwtHandling]()
}

func (e environment) GetDBAccess() (db.DBAccess, error) {
	return env.ParseAs[db.DBAccess]()
}
