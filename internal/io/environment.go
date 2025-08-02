package io

import (
	"github.com/caarlos0/env/v11"
	"github.com/motojouya/geezer_auth/internal/core/config"
	"github.com/motojouya/geezer_auth/pkg/core/jwt"
)

type Environment interface {
	GetJwtHandle() (jwt.JwtHandle, error)
	GetDBAccess() (config.DBAccess, error)
}

type environment struct{}

func CreateEnvironment() Environment {
	return &environment{}
}

type JwtHandleGetter interface {
	GetJwtHandle() (jwt.JwtHandle, error)
}

func (e environment) GetJwtHandle() (jwt.JwtHandle, error) {
	return env.ParseAs[jwt.JwtHandle]()
}

type DBAccessGetter interface {
	GetDBAccess() (config.DBAccess, error)
}

func (e environment) GetDBAccess() (config.DBAccess, error) {
	return env.ParseAs[config.DBAccess]()
}
