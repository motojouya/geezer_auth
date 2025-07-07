package service

import (
	"github.com/motojouya/geezer_auth/internal/io"
	core "github.com/motojouya/geezer_auth/internal/core/db"
	"github.com/motojouya/geezer_auth/internal/db"
)

type DatabaseLoader interface {
	LoadDatabase(e io.Environment) (db.ORP, error)
}

type databaseLoaderImpl struct{}

var dbAccess db.DBAccess

func (imple databaseLoaderImpl) LoadDatabase(e io.Environment) (db.ORP, error) {
	// access 情報はcacheするが、connectionはcacheしない
	if dbAccess != nil {
		var dbAccessData, err = e.GetDBAccess()
		if err != nil {
			return dbAccessData, err
		}
		dbAccess = dbAccessData
	}

	var connection, err = core.CreateConnection(dbAccess)
	if err != nil {
		return connection, err
	}

	var database, err = db.CreateDatabase(connection)
	if err != nil {
		return database, err
	}

	return database, nil
}
