package command

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"time"
)

type ExpireRefreshTokenQuery interface {
	ExpireRefreshToken(userPersistKey uint, now time.Time) error
}

func ExpireRefreshToken(executer gorp.SqlExecutor, userPersistKey uint, now time.Time) error {
	var sql, args, sqlErr = utility.Dialect.Update("user_refresh_token").Set(goqu.Record{"expire_date": now}).Where(
		goqu.C("user_persist_key").Eq(userPersistKey),
		goqu.C("expire_date").IsNull(),
	).Prepared(true).ToSQL()
	if sqlErr != nil {
		return sqlErr
	}

	var _, execErr = executer.Exec(sql, args...)
	if execErr != nil {
		return execErr
	}

	return nil
}
