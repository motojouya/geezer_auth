package service

import (
	"github.com/motojouya/geezer_auth/internal/core/config"
	"github.com/motojouya/geezer_auth/internal/db"
	"github.com/motojouya/geezer_auth/internal/io"
)

type DatabaseLoader interface {
	LoadDatabase(e io.Environment) (db.ORP, error)
}

type databaseLoaderImpl struct{}

var dbAccess *config.DBAccess

func (imple databaseLoaderImpl) LoadDatabase(e io.Environment) (db.ORP, error) {
	// access 情報はcacheするが、connectionはcacheしない
	if dbAccess == nil {
		var dbAccessData, err = e.GetDBAccess()
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
