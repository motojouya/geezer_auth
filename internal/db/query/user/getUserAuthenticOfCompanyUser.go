package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"time"
)

type GetUserAuthenticOfCompanyUserQuery interface {
	GetUserAuthenticOfCompanyUser(companyIdentifier string, userIdentifier string, now time.Time) (*transfer.UserAuthentic, error)
}

func GetUserAuthenticOfCompanyUser(executer gorp.SqlExecutor, companyIdentifier string, userIdentifier string, now time.Time) (*transfer.UserAuthentic, error) {
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
			goqu.I("c.identifier").Eq(companyIdentifier),
		),
	).Where(
		goqu.I("u.identifier").Eq(userIdentifier),
	).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var keys = map[string]string{"company_identifier": companyIdentifier, "user_identifier": userIdentifier}
	var ua, execErr = utility.SelectSingle[transfer.UserAuthentic](executer, "user", keys, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	if ua == nil {
		return nil, nil
	}

	var ucrs, getUserCompanyRolesErr = GetUserCompanyRole(executer, []string{userIdentifier}, now)
	if getUserCompanyRolesErr != nil {
		return nil, getUserCompanyRolesErr
	}

	ua.UserCompanyRole = ucrs
	return ua, nil
}
