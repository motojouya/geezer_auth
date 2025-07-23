package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"time"
)

type GetUserCompanyRoleQuery interface {
	GetUserCompanyRole(identifiers []string, now time.Time) ([]transfer.UserCompanyRoleFull, error)
}

func GetUserCompanyRole(executer gorp.SqlExecutor, identifiers []string, now time.Time) ([]transfer.UserCompanyRoleFull, error) {
	var sql, args, sqlErr = transfer.SelectUserCompanyRole.Where(
		goqu.I("u.identifier").In(identifiers),
		goqu.Or(
			goqu.I("ucr.expire_date").Gte(now),
			goqu.I("ucr.expire_date").IsNull(),
		),
	).Order(
		goqu.I("ucr.user_persist_key").Asc(),
		goqu.I("ucr.company_persist_key").Asc(),
		goqu.I("ucr.role_label").Asc(),
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
