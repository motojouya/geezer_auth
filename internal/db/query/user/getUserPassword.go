package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"time"
)

type GetUserPasswordQuery interface {
	GetUserPassword(identifier string, now time.Time) (*transfer.UserPasswordFull, error)
}

func GetUserPassword(executer gorp.SqlExecutor, identifier string, now time.Time) (*transfer.UserPasswordFull, error) {
	var sql, args, sqlErr = transfer.SelectUserPassword.Where(
		goqu.C("u.identifier").Eq(identifier),
		goqu.Or(
			goqu.C("up.expire_date").Gte(now),
			goqu.C("up.expire_date").IsNull(),
		),
	).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var up, execErr = utility.SelectSingle[transfer.UserPasswordFull](executer, "user_password", map[string]string{"identifier": identifier}, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	return up, nil
}

// TODO verify_date, expire_dateを見るかいなか。verification最中のemailをどう許容するのか。一つしかない状況を作るほうが健全か。
// 当該のemailを使っていいか否かだけど、これはexpire_date is not nullならいいか。間違えて人のやつを試してみることはあり得る。
// verification中のものが複数ある場合、どれにすべきか判断が面倒になるので。なので、新しくemail登録がなされた場合は、verification中のものはexpire_dateを設定するようにする。
