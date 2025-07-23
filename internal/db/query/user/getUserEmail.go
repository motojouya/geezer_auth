package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
)

type GetUserEmailQuery interface {
	GetUserEmail(email string) ([]transfer.UserEmailFull, error)
}

func GetUserEmail(executer gorp.SqlExecutor, email string) ([]transfer.UserEmailFull, error) {
	var sql, args, sqlErr = transfer.SelectUserEmail.Where(
		goqu.I("ue.email").Eq(email),
		goqu.I("ue.expire_date").IsNull(),
	).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var ues []transfer.UserEmailFull
	var _, execErr = executer.Select(&ues, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	return ues, nil
}

// TODO verify_date, expire_dateを見るかいなか。verification最中のemailをどう許容するのか。一つしかない状況を作るほうが健全か。
// 当該のemailを使っていいか否かだけど、これはexpire_date is not nullならいいか。間違えて人のやつを試してみることはあり得る。
// verification中のものが複数ある場合、どれにすべきか判断が面倒になるので。なので、新しくemail登録がなされた場合は、verification中のものはexpire_dateを設定するようにする。
