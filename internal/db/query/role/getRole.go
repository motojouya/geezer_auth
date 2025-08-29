package role

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"github.com/motojouya/geezer_auth/internal/db/utility"
)

type GetRoleQuery interface {
	GetRole(label string) (*transfer.Role, error)
}

func GetRole(executer gorp.SqlExecutor, label string) (*transfer.Role, error) {
	var sql, args, sqlErr = transfer.SelectRole.Where(goqu.I("r.label").Eq(label)).Order(goqu.I("r.label").Asc()).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var role, execErr = utility.SelectSingle[transfer.Role](executer, "role", map[string]string{"label": label}, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	return role, nil
}
