package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type DBAccess struct {
	Driver string `env:"DB_DRIVER,notEmpty"`
	Name   string `env:"DB_NAME,notEmpty"`
	Host   string `env:"DB_HOST,notEmpty"`
	Port   uint   `env:"DB_PORT,notEmpty"`
	User   string `env:"DB_USER,notEmpty"`
	Pass   string `env:"DB_PASSWORD,notEmpty"`
}

func NewDBAccess(
	driver string,
	name string,
	host string,
	port uint,
	user string,
	pass string,
) DBAccess {
	return DBAccess{
		Driver: driver,
		Name:   name,
		Host:   host,
		Port:   port,
		User:   user,
		Pass:   pass,
	}
}

func (dbAccess DBAccess) CreateConnection() (*sql.DB, error) {
	return sql.Open(dbAccess.Driver, fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		dbAccess.User,
		dbAccess.Pass,
		dbAccess.Host,
		dbAccess.Port,
		dbAccess.Name,
	))
}
