package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"time"
)

type GetUserRefreshTokenQuery interface {
	GetUserRefreshToken(identifier string, now time.Time) (*transfer.UserRefreshTokenFull, error)
}

func GetUserRefreshToken(executer gorp.SqlExecutor, identifier string, now time.Time) (*transfer.UserRefreshTokenFull, error) {
	var sql, args, sqlErr = transfer.SelectUserRefreshToken.Where(
		goqu.I("u.identifier").Eq(identifier),
		goqu.I("urt.expire_date").Gte(now),
	).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var urt, execErr = utility.SelectSingle[transfer.UserRefreshTokenFull](executer, "user_refresh_token", map[string]string{"identifier": identifier}, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	return urt, nil
}
