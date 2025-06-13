package io

import (
	"github.com/caarlos0/env/v11"
	"github.com/motojouya/geezer_auth/pkg/core/jwt"
)

type Environment interface {
	GetJwtHandling() (jwt.JwtHandling, error)
}

type environment struct{}

func CreateEnvironment() Environment {
	return &environment{}
}

func (e environment) GetJwtHandling() (jwt.JwtHandling, error) {
	return env.ParseAs[jwt.JwtHandling]()
}
