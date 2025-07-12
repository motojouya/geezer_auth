package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"time"
)

type GetUserEmailOfTokenQuery interface {
	GetUserEmailOfToken(identifier string, email string, verifyToken string, now time.Time) ([]transfer.UserEmail, error)
}

func GetUserEmailOfToken(executer gorp.SqlExecutor, identifier string, email string, verifyToken string, now time.Time) ([]transfer.UserEmail, error) {
	var sql, args, sqlErr = transfer.SelectUserEmail.Where(
		goqu.C("u.identifier").Eq(identifier),
		goqu.C("ue.email").Eq(email),
		goqu.C("ue.verify_token").Eq(verifyToken),
		// goqu.C("ue.verify_date").IsNull(), すでにverifiedならば、登録されましたと返せばいいだけ
		goqu.Or(
			goqu.C("ue.expire_date").Gte(now),
			goqu.C("ue.expire_date").IsNull(),
		),
	).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var ues []transfer.UserEmail
	var _, execErr = executer.Select(&ues, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	return ues, nil
}

// TODO verify_date, expire_dateを見るかいなか。verification最中のemailをどう許容するのか。一つしかない状況を作るほうが健全か。
// 当該のemailを使っていいか否かだけど、これはexpire_date is not nullならいいか。間違えて人のやつを試してみることはあり得る。
// verification中のものが複数ある場合、どれにすべきか判断が面倒になるので。なので、新しくemail登録がなされた場合は、verification中のものはexpire_dateを設定するようにする。
