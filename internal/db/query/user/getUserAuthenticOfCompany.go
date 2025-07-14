package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"time"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"github.com/motojouya/geezer_auth/internal/core/essence"
)

type GetUserAuthenticOfCompanyQuery interface {
	GetUserAuthenticOfCompany(identifier string, now time.Time) ([]transfer.UserAuthentic, error)
}

func GetUserAuthenticOfCompany(executer gorp.SqlExecutor, identifier string, now time.Time) ([]transfer.UserAuthentic, error) {
	var sql, args, sqlErr = transfer.SelectUserAuthentic.InnerJoin(
		utility.Dialect.From("user_company_role").As("ucr").Select(goqu.Select("ucr.company_persist_key").Distinct().As("company_persist_key")),
		goqu.On(
			goqu.C("u.persist_key").Eq("ucr.user_persist_key"),
			goqu.Or(
				goqu.C("ucr.expire_date").Gte(now),
				goqu.C("ucr.expire_date").IsNull(),
			),
		),
	).InnerJoin(
		utility.Dialect.From("company").As("c"),
		goqu.On(
			goqu.C("c.persist_key").Eq("ucr.company_persist_key"),
			goqu.C("c.identifier").Eq(identifier),
		),
	).Prepared(true).ToSQL()
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
