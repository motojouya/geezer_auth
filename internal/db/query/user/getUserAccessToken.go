package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"time"
)

type GetUserAccessTokenQuery interface {
	GetUserAccessToken(identifier string, now time.Time) ([]transfer.UserAccessToken, error)
}

func GetUserAccessToken(executer gorp.SqlExecutor, identifier string, now time.Time) ([]transfer.UserAccessToken, error) {
	var sql, args, sqlErr = transfer.SelectUserAccessToken.Where(
		goqu.C("u.identifier").Eq(identifier),
		goqu.C("uat.source_update_date").Eq("u.update_date"),
		goqu.C("uat.expire_date").Gte(now),
	).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var uats []transfer.UserAccessToken
	var _, execErr = executer.Select(&uats, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	return uats, nil
}
