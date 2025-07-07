package service

import (
	"github.com/motojouya/geezer_auth/internal/core/authorization"
	"github.com/motojouya/geezer_auth/internal/core/role"
	text "github.com/motojouya/geezer_auth/pkg/core/text"
	user "github.com/motojouya/geezer_auth/pkg/core/user"
)

type AuthorizationLoader interface {
	LoadAuthorization() Authorizer
}

type authorizationLoaderImpl struct{}

type Authorizer interface {
	Authorize(require role.RequirePermission, authentic *user.Authentic) error
}

var authorizationSingleton *authorization.Authorization

// TODO DBアクセスしてロードするが、まだDB実装していない
func (imple authorizationLoaderImpl) LoadAuthorization() Authorizer {
	if authorizationSingleton != nil {
		return authorizationSingleton
	}

	var EmployeeLabel, employeeLabelErr = text.NewLabel("EMPLOYEE")
	if employeeLabelErr != nil {
		panic(employeeLabelErr)
	}
	var EmployeePermission = role.NewRolePermission(EmployeeLabel, true, true, false, false, 5)

	var ManagerLabel, managerLabelErr = text.NewLabel("MANAGER")
	if managerLabelErr != nil {
		panic(managerLabelErr)
	}
	var ManagerPermission = role.NewRolePermission(ManagerLabel, true, true, true, true, 9)

	var permissions = []role.RolePermission{
		EmployeePermission,
		ManagerPermission,
	}

	authorizationSingleton = authorization.CreateAuthorization(permissions)
	return authorizationSingleton
}
