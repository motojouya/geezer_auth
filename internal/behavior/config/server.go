package config

import (
	"github.com/motojouya/geezer_auth/internal/local"
	shelterConfig "github.com/motojouya/geezer_auth/internal/shelter/config"
)

type ServerGetter interface {
	GetServer() (shelterConfig.Server, error)
}

type ServerGet struct {
	env local.ServerGetter
}

func NewServerGet(env local.ServerGetter) *ServerGet {
	return &ServerGet{
		env: env,
	}
}

var serverConf *shelterConfig.Server

func (getter *ServerGet) GetServer() (*shelterConfig.Server, error) {
	if serverConf == nil {
		var serverConfObj, err = getter.env.GetServer()
		if err != nil {
			return nil, err
		}

		serverConf = &serverConfObj
	}

	return serverConf, nil
}
