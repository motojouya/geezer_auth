package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/core/essence"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"time"
)

type GetUserAuthenticOfCompanyQuery interface {
	GetUserAuthenticOfCompany(identifier string, now time.Time) ([]transfer.UserAuthentic, error)
}

func GetUserAuthenticOfCompany(executer gorp.SqlExecutor, identifier string, now time.Time) ([]transfer.UserAuthentic, error) {
	var sql, args, sqlErr = transfer.SelectUserAuthentic.InnerJoin(
		utility.Dialect.From(goqu.T("user_company_role").As("ucr")).Where(
			goqu.Or(
				goqu.I("ucr.expire_date").Gte(now),
				goqu.I("ucr.expire_date").IsNull(),
			),
		).Select(
			goqu.I("ucr.user_persist_key").As("user_persist_key"),
			goqu.I("ucr.company_persist_key").As("company_persist_key"),
		).Distinct().As("cref"),
		goqu.On(
			goqu.I("u.persist_key").Eq(goqu.I("cref.user_persist_key")),
		),
	).InnerJoin(
		goqu.T("company").As("c"),
		goqu.On(
			goqu.I("c.persist_key").Eq(goqu.I("cref.company_persist_key")),
			goqu.I("c.identifier").Eq(identifier),
		),
	).Order(goqu.I("u.persist_key").Asc()).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var uas []transfer.UserAuthentic
	var _, execErr = executer.Select(&uas, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	if len(uas) == 0 {
		return []transfer.UserAuthentic{}, nil
	}

	var uaPtrs = essence.ToPtr(uas)
	var userIdentifiers = essence.Map(uaPtrs, transfer.UserIdentifierUserAuthentic)

	var ucrs, getUserCompanyRolesErr = GetUserCompanyRole(executer, userIdentifiers, now)
	if getUserCompanyRolesErr != nil {
		return nil, getUserCompanyRolesErr
	}

	var ptrs = essence.Relate(uaPtrs, ucrs, transfer.RelateUserCompanyRole)
	return essence.ToVal(ptrs), nil
}
