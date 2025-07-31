package authorization

import (
	"github.com/motojouya/geezer_auth/internal/core/authorization"
	"github.com/motojouya/geezer_auth/internal/core/role"
	roleQuery "github.com/motojouya/geezer_auth/internal/db/query/role"
	user "github.com/motojouya/geezer_auth/pkg/core/user"
)

type AuthorizationLoader struct {
	db roleQuery.GetRolePermissionQuery
}

func NewAuthorizationLoader(db roleQuery.GetRolePermissionQuery) *AuthorizationLoader {
	return &AuthorizationLoader{db: db}
}

type Authorizer interface {
	Authorize(require role.RequirePermission, authentic *user.Authentic) error
}

var authorizationSingleton *authorization.Authorization

func (loader AuthorizationLoader) LoadAuthorization() (Authorizer, error) {
	if authorizationSingleton != nil {
		return authorizationSingleton, nil
	}

	dbRolePermissions, err := loader.db.GetRolePermission()
	if err != nil {
		return nil, err
	}

	permissions := make([]role.RolePermission, len(dbRolePermissions))
	for _, dbRolePermission := range dbRolePermissions {
		rolePermission, err := dbRolePermission.ToCoreRolePermission()
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, rolePermission)
	}

	// var EmployeeLabel, employeeLabelErr = text.NewLabel("EMPLOYEE")
	// if employeeLabelErr != nil {
	// 	panic(employeeLabelErr)
	// }
	// var EmployeePermission = role.NewRolePermission(EmployeeLabel, true, true, false, false, 5)

	// var ManagerLabel, managerLabelErr = text.NewLabel("MANAGER")
	// if managerLabelErr != nil {
	// 	panic(managerLabelErr)
	// }
	// var ManagerPermission = role.NewRolePermission(ManagerLabel, true, true, true, true, 9)

	// var permissions = []role.RolePermission{
	// 	EmployeePermission,
	// 	ManagerPermission,
	// }

	authorizationSingleton = authorization.CreateAuthorization(permissions)
	return authorizationSingleton, nil
}
