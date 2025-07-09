package query

import (
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
)

type GetRoleQuery interface {
	GetRole() ([]role.Role, error)
}

func GetRole(executer gorp.SqlExecutor) ([]role.Role, error) {
	var sql, _, sqlErr = role.SelectRole.ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var roles []role.Role
	var _, execErr = executer.Select(&roles, sql)
	if execErr != nil {
		return nil, execErr
	}

	return roles, nil
}
