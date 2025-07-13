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
		goqu.C("u.identifier").Eq(identifier),
		goqu.C("urt.expire_date").Gte(now),
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

// TODO verify_date, expire_dateを見るかいなか。verification最中のemailをどう許容するのか。一つしかない状況を作るほうが健全か。
// 当該のemailを使っていいか否かだけど、これはexpire_date is not nullならいいか。間違えて人のやつを試してみることはあり得る。
// verification中のものが複数ある場合、どれにすべきか判断が面倒になるので。なので、新しくemail登録がなされた場合は、verification中のものはexpire_dateを設定するようにする。
