package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Connection interface {
	Close() error
}

func CreateConnection(dbAccess DBAccess) (*sql.DB, error) {
	return sql.Open(dbAccess.Driver, fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		dbAccess.User,
		dbAccess.Pass,
		dbAccess.Host,
		dbAccess.Port,
		dbAccess.Name,
	))
}
