package config

import (
	coreConfig "github.com/motojouya/geezer_auth/internal/core/config"
	"github.com/motojouya/geezer_auth/internal/db"
)

type DBAccessGetter interface {
	GetDBAccess() (coreConfig.DBAccess, error)
}

type DatabaseLoader struct {
	env DBAccessGetter
}

func NewDatabaseLoader(env DBAccessGetter) *DatabaseLoader {
	return &DatabaseLoader{
		env: env,
	}
}

var dbAccess *coreConfig.DBAccess

func (loader DatabaseLoader) LoadDatabase() (db.ORP, error) {
	// access 情報はcacheするが、connectionはcacheしない
	if dbAccess == nil {
		var dbAccessData, err = loader.env.GetDBAccess()
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
