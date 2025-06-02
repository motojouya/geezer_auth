package service

import (
	"github.com/motojouya/geezer_auth/internal/core/authorization"
	"github.com/motojouya/geezer_auth/internal/core/role"
	text "github.com/motojouya/geezer_auth/pkg/core/text"
)

type Authorizer interface {
	Authorize(require RequirePermission, authentic *user.Authentic) error
}

// TODO DBアクセスしてロードするが、まだDB実装していない
func LoadAuthorization() Authorizer {
	var EmployeeLabel = text.NewLabel("EMPLOYEE")
	var EmployeePermission = role.NewRolePermission(EmployeeLabel, true, true, false, false, 5)

	var ManagerLabel = text.NewLabel("MANAGER")
	var ManagerPermission = role.NewRolePermission(ManagerLabel, true, true, true, true, 9)

	var permissions = []role.RolePermission{
		EmployeePermission,
		ManagerPermission,
	}

	return authorization.CreateAuthorization(permissions)
}
