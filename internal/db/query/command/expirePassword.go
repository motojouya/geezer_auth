package command

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"time"
)

type ExpirePasswordQuery interface {
	ExpirePassword(userPersistKey uint, now time.Time) error
}

func ExpirePassword(executer gorp.SqlExecutor, userPersistKey uint, now time.Time) error {
	var sql, args, sqlErr = utility.Dialect.Update("user_password").Set(goqu.Record{"expire_date": now}).Where(
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
