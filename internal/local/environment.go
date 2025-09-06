package local

import (
	"github.com/caarlos0/env/v11"
	"github.com/motojouya/geezer_auth/internal/shelter/config"
	"github.com/motojouya/geezer_auth/pkg/shelter/jwt"
)

type Environmenter interface {
	JwtHandleGetter
	DBAccessGetter
}

type Environment struct{}

func CreateEnvironment() *Environment {
	return &Environment{}
}

type JwtHandleGetter interface {
	GetJwtHandle() (jwt.JwtHandle, error)
}

func (e Environment) GetJwtHandle() (jwt.JwtHandle, error) {
	return env.ParseAs[jwt.JwtHandle]()
}

type DBAccessGetter interface {
	GetDBAccess() (config.DBAccess, error)
}

func (e Environment) GetDBAccess() (config.DBAccess, error) {
	return env.ParseAs[config.DBAccess]()
}

type ServerGetter interface {
	GetServer() (config.Server, error)
}

func (e Environment) GetServer() (config.Server, error) {
	return env.ParseAs[config.Server]()
}
