package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/motojouya/geezer_auth/internal/db/utility"
)

type GetUserQuery interface {
	GetUser(identifier string) (*transfer.User, error)
}

func GetUser(executer gorp.SqlExecutor, identifier string) (*transfer.User, error) {
	var sql, args, sqlErr = transfer.SelectUser.Where(goqu.I("u.identifier").Eq(identifier)).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	return utility.SelectSingle[transfer.User](executer, "users", map[string]string{"identifier": identifier}, sql, args...)
}
