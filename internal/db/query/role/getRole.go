package role

import (
	"github.com/go-gorp/gorp"
	"github.com/doug-martin/goqu/v9"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/role"
)

type GetRoleQuery interface {
	GetRole() ([]transfer.Role, error)
}

func GetRole(executer gorp.SqlExecutor) ([]transfer.Role, error) {
	var sql, _, sqlErr = transfer.SelectRole.Order(goqu.I("r.label").Asc()).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var roles []transfer.Role
	var _, execErr = executer.Select(&roles, sql)
	if execErr != nil {
		return nil, execErr
	}

	return roles, nil
}
