package db

import (
	"database/sql"
	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
	"github.com/motojouya/geezer_auth/internal/core/essence"
	"github.com/motojouya/geezer_auth/internal/db/transfer/company"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
)

// FIXME Prepare関数いる？
// TODO migrationは、このモジュールではなく、dbモジュールを起動時に呼び出して行うので、webプロセスではしない
// TODO defer dbMap.Db.Close() は、内部にConnectionを持っている場合、自動で呼び出せるように工夫する。これはwebのmiddlewareで行う

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
	Command
}

func CreateDatabase(connection *sql.DB) ORP {
	var dbMap = &gorp.DbMap{Db: connection, Dialect: gorp.PostgresDialect{}}
	registerTable(dbMap)

	return &ORPImpl{
		DbMap: dbMap,
	}
}

func registerTable(dbMap *gorp.DbMap) {
	company.AddCompanyTable(dbMap)
	company.AddCompanyInviteTable(dbMap)
	role.AddRoleTable(dbMap)
	role.AddRolePermissionTable(dbMap)
	user.AddUserTable(dbMap)
	//user.AddUserAccessTokenTable(dbMap)
	//user.AddUserCompanyRoleTable(dbMap)
	//user.AddUserEmailTable(dbMap)
	//user.AddUserPasswordTable(dbMap)
	//user.AddUserRefreshTokenTable(dbMap)
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
