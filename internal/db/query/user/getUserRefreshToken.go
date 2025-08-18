package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"time"
)

type GetUserRefreshTokenQuery interface {
	GetUserRefreshToken(token string, now time.Time) (*transfer.UserAuthentic, error)
}

func GetUserRefreshToken(executer gorp.SqlExecutor, token string, now time.Time) (*transfer.UserAuthentic, error) {
	var sql, args, sqlErr = transfer.SelectUserAuthentic.InnerJoin(
		goqu.T("user_refresh_token").As("urt"),
		goqu.On(
			goqu.I("u.persist_key").Eq(goqu.I("urt.user_persist_key")),
			goqu.I("urt.refresh_token").Eq(token)
			goqu.I("urt.expire_date").Gte(now),
		),
	).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var ua, execErr = utility.SelectSingle[transfer.UserAuthentic](executer, "user", map[string]string{"refreshToken": token}, sql, args...)
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
