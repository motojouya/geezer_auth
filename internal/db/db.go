package db

import (
	core "github.com/motojouya/geezer_auth/internal/core/db"
	"database/sql"
	"fmt"
	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

// FIXME Prepare関数いる？
// TODO migrationは、このモジュールではなく、dbモジュールを起動時に呼び出して行うので、webプロセスではしない
// TODO defer dbMap.Db.Close() は、内部にConnectionを持っている場合、自動で呼び出せるように工夫する。これはwebのmiddlewareで行う

type ORP interface {
	gorp.SqlExecutor
	Begin() (ORPTransaction, error)
	core.Connection
}

type ORPTransaction interface {
	gorp.SqlExecutor
	Commit() error
	Rollback() error
}

func CreateDatabase(connection *sql.DB) ORP {
	var dbMap = &gorp.DbMap{Db: connection, Dialect: gorp.PostgresDialect{}}
	return &ORPImpl{
		DbMap: dbMap,
	}
}

type ORPImpl struct {
	*gorp.DbMap
}

type ORPTransactionImpl struct {
	*gorp.Transaction
}

func (orp ORPImpl) Close() error {
	return orp.DbMap.Db.Close()
}

func (orp ORPImpl) Begin() (ORPTransaction, error) {
	var transaction, err = orp.DbMap.Begin()
	if err != nil {
		return nil, err
	}

	return &ORPTransactionImpl{
		Transaction: transaction,
	}, nil
}
