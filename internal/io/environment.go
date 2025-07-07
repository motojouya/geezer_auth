package io

import (
	"github.com/caarlos0/env/v11"
	"github.com/motojouya/geezer_auth/pkg/core/jwt"
	"github.com/motojouya/geezer_auth/internal/core/db"
)

type Environment interface {
	GetJwtHandling() (jwt.JwtHandling, error)
	GetDBAccess() (db.DbAccess, error)
}

type environment struct{}

func CreateEnvironment() Environment {
	return &environment{}
}

func (e environment) GetJwtHandling() (jwt.JwtHandling, error) {
	return env.ParseAs[jwt.JwtHandling]()
}

func (e environment) GetDBAccess() (db.DbAccess, error) {
	return env.ParseAs[db.DbAccess]()
}
