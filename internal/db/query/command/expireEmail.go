package command

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	"time"
	"github.com/motojouya/geezer_auth/internal/db/utility"
)

type ExpireEmailQuery interface {
	ExpireEmail(userPersistKey uint, persistKey uint, ignoreVerified bool, now time.Time) error
}

func ExpireEmail(executer gorp.SqlExecutor, userPersistKey uint, persistKey uint, ignoreVerified bool, now time.Time) error {
	var query = utility.Dialect.Update("user_email").Set(goqu.Record{"expire_date": now}).Where(
		goqu.C("user_persist_key").Eq(userPersistKey),
		goqu.C("persist_key").Neq(persistKey),
		goqu.C("expire_date").IsNull(),
	)
	if ignoreVerified {
		query = query.Where(goqu.C("ue.verify_date").IsNull())
	}

	var sql, args, sqlErr = query.Prepared(true).ToSQL()
	if sqlErr != nil {
		return sqlErr
	}

	var _, execErr = executer.Exec(sql, args...)
	if execErr != nil {
		return execErr
	}

	return nil
}
