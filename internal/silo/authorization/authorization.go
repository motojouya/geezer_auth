package authorization

import (
	"github.com/motojouya/geezer_auth/internal/core/authorization"
	"github.com/motojouya/geezer_auth/internal/core/role"
	roleQuery "github.com/motojouya/geezer_auth/internal/db/query/role"
)

type AuthorizationGetter interface {
	GetAuthorization() (*authorization.Authorization, error)
}

type AuthorizationGet struct {
	db roleQuery.GetRolePermissionQuery
}

func NewAuthorizationGet(db roleQuery.GetRolePermissionQuery) *AuthorizationGet {
	return &AuthorizationGet{db: db}
}

var authorizationSingleton *authorization.Authorization

func (getter AuthorizationGet) GetAuthorization() (*authorization.Authorization, error) {
	if authorizationSingleton != nil {
		return authorizationSingleton, nil
	}

	dbRolePermissions, err := getter.db.GetRolePermission()
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

	authorizationSingleton = authorization.CreateAuthorization(permissions)
	return authorizationSingleton, nil
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
