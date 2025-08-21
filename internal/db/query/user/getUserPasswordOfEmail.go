package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/motojouya/geezer_auth/internal/db/utility"
)

type GetUserPasswordOfEmailQuery interface {
	GetUserPasswordOfEmail(email string) (*transfer.UserPasswordFull, error)
}

func GetUserPasswordOfEmail(executer gorp.SqlExecutor, email string) (*transfer.UserPasswordFull, error) {
	var sql, args, sqlErr = transfer.SelectUserPassword.Where(
		goqu.I("u.email_identifier").Eq(email),
		goqu.I("up.expire_date").IsNull(),
	).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var up, execErr = utility.SelectSingle[transfer.UserPasswordFull](executer, "user_password", map[string]string{"email_identifier": email}, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	return up, nil
}
