package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"time"
	"github.com/motojouya/geezer_auth/internal/db/utility"
)

type GetUserAuthenticQuery interface {
	GetUserAuthentic(identifier string, now time.Time) (*transfer.UserAuthentic, error)
}

func GetUserAuthentic(executer gorp.SqlExecutor, identifier string, now time.Time) (*transfer.UserAuthentic, error) {
	var sql, args, sqlErr = transfer.SelectUserAuthentic.Where(goqu.C("u.identifier").Eq(identifier)).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var ua, execErr = utility.SelectSingle[transfer.UserAuthentic](executer, "user", map[string]string{"identifier":identifier}, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	if ua == nil {
		return nil, nil
	}

	var ucrs, getUserCompanyRolesErr = GetUserCompanyRole(executer, []string{identifier}, now)
	if getUserCompanyRolesErr != nil {
		return nil, getUserCompanyRolesErr
	}

	ua.UserCompanyRole = ucrs
	return ua, nil
}
