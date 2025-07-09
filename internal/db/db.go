package db

import (
	"database/sql"
	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
	"github.com/motojouya/geezer_auth/internal/core/essence"
)

// FIXME Prepare関数いる？
// TODO migrationは、このモジュールではなく、dbモジュールを起動時に呼び出して行うので、webプロセスではしない
// TODO defer dbMap.Db.Close() は、内部にConnectionを持っている場合、自動で呼び出せるように工夫する。これはwebのmiddlewareで行う
// TODO primary key の設定いるかな？テーブル名は合わせてるので、不要だが、autoincrementはどうか。必要なら関数をtransferに用意して、こっちで呼び出す感じ
// t1 := dbmap.AddTableWithName(Invoice{}, "invoice_test").SetKeys(true, "Id")

type ORP interface {
	gorp.SqlExecutor
	Begin() (ORPTransaction, error)
	essence.Closable
	Query
}

type ORPTransaction interface {
	gorp.SqlExecutor
	Commit() error
	Rollback() error
	Query
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
