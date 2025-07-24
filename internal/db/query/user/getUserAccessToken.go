package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"time"
)

type GetUserAccessTokenQuery interface {
	GetUserAccessToken(identifier string, now time.Time) ([]transfer.UserAccessTokenFull, error)
}

/*
 * 以下でも良かったが、1,2病の間に登録される可能性は低いので、複数レコードといっても数レコードの想定。
 * `goqu.C("uat.source_update_date").Eq("u.update_date"),`
 * それよりは、access tokenが取得できない状況のほうが問題なので、betweenで広く取得する。
 */
func GetUserAccessToken(executer gorp.SqlExecutor, identifier string, now time.Time) ([]transfer.UserAccessTokenFull, error) {
	var sql, args, sqlErr = transfer.SelectUserAccessToken.Where(
		goqu.I("u.identifier").Eq(identifier),
		goqu.I("uat.source_update_date").Between(goqu.Range(goqu.L("u.update_date + '-1 second'"), goqu.L("u.update_date + '1 second'"))),
		goqu.I("uat.expire_date").Gte(now),
	).Order(
		goqu.I("uat.register_date").Desc(),
	).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var uats []transfer.UserAccessTokenFull
	var _, execErr = executer.Select(&uats, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	return uats, nil
}
