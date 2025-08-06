package config

import (
	"github.com/motojouya/geezer_auth/internal/db"
	"github.com/motojouya/geezer_auth/internal/local"
	shelterConfig "github.com/motojouya/geezer_auth/internal/shelter/config"
)

type DatabaseGetter interface {
	GetDatabase() (*db.ORP, error)
}

type DatabaseGet struct {
	env local.DBAccessGetter
}

func NewDatabaseGet(env local.DBAccessGetter) *DatabaseGet {
	return &DatabaseGet{
		env: env,
	}
}

var dbAccess *shelterConfig.DBAccess

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
