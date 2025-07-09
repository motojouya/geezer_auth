package query

import (
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
)

type GetRolePermissionQuery interface {
	GetRolePermission() ([]role.RolePermission, error)
}

func GetRolePermission(executer gorp.SqlExecutor) ([]role.RolePermission, error) {
	var sql, _, sqlErr = role.SelectRolePermission.ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var rolePermissions []role.RolePermission
	var _, execErr = executer.Select(&rolePermissions, sql)
	if execErr != nil {
		return nil, execErr
	}

	return rolePermissions, nil
}
