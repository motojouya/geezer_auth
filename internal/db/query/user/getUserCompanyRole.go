package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"time"
)

type GetUserCompanyRoleQuery interface {
	GetUserCompanyRole(identifier string, now time.Time) ([]transfer.UserCompanyRoleFull, error)
}

func GetUserCompanyRole(executer gorp.SqlExecutor, identifier string, now time.Time) ([]transfer.UserCompanyRoleFull, error) {
	var sql, args, sqlErr = transfer.SelectUserCompanyRole.Where(
		goqu.C("u.identifier").Eq(identifier),
		goqu.Or(
			goqu.C("ucr.expire_date").Gte(now),
			goqu.C("ucr.expire_date").IsNull(),
		),
	).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var ucrs []transfer.UserCompanyRoleFull
	var _, execErr = executer.Select(&ucrs, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	return ucrs, nil
}
