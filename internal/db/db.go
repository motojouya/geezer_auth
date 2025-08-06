package db

import (
	"database/sql"
	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
	"github.com/motojouya/geezer_auth/internal/db/transfer/company"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
)

// FIXME Prepare関数いる？
// TODO defer dbMap.Db.Close() は、内部にConnectionを持っている場合、自動で呼び出せるように工夫する。これはwebのmiddlewareで行う

type Transactional interface {
	Begin() error
	Commit() error
	Rollback() error
}

type TransactionalDatabase interface {
	Transactional
	essence.Closable
}

type ORPer interface {
	TransactionalDatabase
	gorp.SqlExecutor
	Query
}

func CreateDatabase(connection *sql.DB) *ORP {
	var dbMap = &gorp.DbMap{Db: connection, Dialect: gorp.PostgresDialect{}}
	registerTable(dbMap)

	return &ORP{
		SqlExecutor: dbMap,
		dbMap:       nil,
	}
}

func registerTable(dbMap *gorp.DbMap) {
	user.AddUserAccessTokenTable(dbMap)
	user.AddUserRefreshTokenTable(dbMap)
	user.AddUserPasswordTable(dbMap)
	user.AddUserCompanyRoleTable(dbMap)
	company.AddCompanyInviteTable(dbMap)
	user.AddUserEmailTable(dbMap)
	company.AddCompanyTable(dbMap)
	user.AddUserTable(dbMap)
	role.AddRolePermissionTable(dbMap)
	role.AddRoleTable(dbMap)
}

// dbMapはtransactionを開始した際に、退避するためのフィールドなので、transactionが開始されていない場合はnil。
type ORP struct {
	gorp.SqlExecutor
	dbMap *gorp.DbMap
}

func (orp ORP) Close() error {
	var insideTransaction = false

	if orp.dbMap != nil {
		insideTransaction = true
	}

	var dbMap, ok = orp.SqlExecutor.(*gorp.DbMap)
	if !ok {
		var err = orp.Rollback()
		if err != nil {
			return essence.CreateInsideTransactionError("transaction is not closed yet. and cannot closed transaction and connection.")
		}

		// rollbackしているので、`gorp.DbMap`になっているはず。失敗しているならいずれにしろcloseできないので、↑のreturnでerrorが返る。
		dbMap, ok = orp.SqlExecutor.(*gorp.DbMap)
		if !ok {
			return essence.CreateInsideTransactionError("transaction is not closed yet. and cannot closed transaction and connection.")
		}
		insideTransaction = true
	}

	var err = dbMap.Db.Close()
	if err != nil {
		return err
	}

	// closeは基本的に強制的に行うが、transactionが開いていた場合は、関数としてはエラーとする。
	if insideTransaction {
		return essence.CreateExitTransactionError("transaction is not closed yet. but closed transaction and connection already.")
	}

	return nil
}

func (orp ORP) Begin() error {
	var dbMap, ok = orp.SqlExecutor.(*gorp.DbMap)
	if !ok || orp.dbMap != nil {
		return essence.CreateInsideTransactionError("transaction is already started")
	}

	var transaction, err = dbMap.Begin()
	if err != nil {
		return err
	}

	orp.SqlExecutor = transaction
	orp.dbMap = dbMap

	return nil
}

func (orp ORP) Commit() error {
	var transaction, ok = orp.SqlExecutor.(*gorp.Transaction)
	if !ok || orp.dbMap == nil {
		return essence.CreateOutsideTransactionError("transaction is not started")
	}

	var err = transaction.Commit()
	if err != nil {
		return err
	}

	orp.SqlExecutor = orp.dbMap
	orp.dbMap = nil

	return nil
}

func (orp ORP) Rollback() error {
	var transaction, ok = orp.SqlExecutor.(*gorp.Transaction)
	if !ok || orp.dbMap == nil {
		return essence.CreateOutsideTransactionError("transaction is not started")
	}

	var err = transaction.Rollback()
	if err != nil {
		return err
	}

	orp.SqlExecutor = orp.dbMap
	orp.dbMap = nil

	return nil
}

func (orp ORP) checkTransaction() error {
	var _, ok = orp.SqlExecutor.(*gorp.Transaction)
	if !ok || orp.dbMap == nil {
		return essence.CreateOutsideTransactionError("transaction is not started")
	}

	return nil
}

func (orp ORP) InsideTransaction() bool {
	return orp.dbMap != nil
}
