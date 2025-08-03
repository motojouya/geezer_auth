package config

import (
	coreConfig "github.com/motojouya/geezer_auth/internal/core/config"
	"github.com/motojouya/geezer_auth/internal/db"
	"github.com/motojouya/geezer_auth/internal/io"
)

type DatabaseGetter interface {
	GetDatabase() (*db.ORP, error)
}

type DatabaseGet struct {
	env io.DBAccessGetter
}

func NewDatabaseGet(env io.DBAccessGetter) *DatabaseGet {
	return &DatabaseGet{
		env: env,
	}
}

var dbAccess *coreConfig.DBAccess

func (getter DatabaseGet) GetDatabase() (*db.ORP, error) {
	// access 情報はcacheするが、connectionはcacheしない
	if dbAccess == nil {
		var dbAccessData, err = getter.env.GetDBAccess()
		if err != nil {
			return nil, err
		}
		dbAccess = &dbAccessData
	}

	var connection, err = dbAccess.CreateConnection()
	if err != nil {
		return nil, err
	}

	var database = db.CreateDatabase(connection)
	if err != nil {
		return database, err
	}

	return database, nil
}
