package role

import (
	"github.com/go-gorp/gorp"
	"github.com/doug-martin/goqu/v9"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/role"
)

type GetRolePermissionQuery interface {
	GetRolePermission() ([]transfer.RolePermission, error)
}

func GetRolePermission(executer gorp.SqlExecutor) ([]transfer.RolePermission, error) {
	var sql, _, sqlErr = transfer.SelectRolePermission.Order(goqu.I("rp.role_label").Asc()).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var rolePermissions []transfer.RolePermission
	var _, execErr = executer.Select(&rolePermissions, sql)
	if execErr != nil {
		return nil, execErr
	}

	return rolePermissions, nil
}
