package authorization

import (
	"github.com/motojouya/geezer_auth/internal/model/role"
	"github.com/motojouya/geezer_auth/pkg/utility"
	text "github.com/motojouya/geezer_auth/pkg/model/text"
	user "github.com/motojouya/geezer_auth/pkg/model/user"
)

type Authorization struct {
	Permisions []role.RolePermission
}

func NewAuthorization(permissions []role.RolePermission) *Authorization {
	return &Authorization{
		Permisions: permissions,
	}
}

type RequirePermission struct {
	SelfEdit      bool
	CompanyAccess bool
	CompanyInvite bool
	CompanyEdit   bool
}

func NewRequirePermission(selfEdit bool, companyAccess bool, companyInvite bool, companyEdit bool) RequirePermission {
	return RequirePermission{
		SelfEdit:      selfEdit,
		CompanyAccess: companyAccess,
		CompanyInvite: companyInvite,
		CompanyEdit:   companyEdit,
	}
}

// TODO DBアクセスしてロードするが、まだDB実装していない
func CreateAuthorization() *Authorization {
	var EmployeeLabel = text.NewLabel("EMPLOYEE")
	var EmployeePermission = role.NewRolePermission(EmployeeLabel, true, true, false, false, 5)

	var ManagerLabel = text.NewLabel("MANAGER")
	var ManagerPermission = role.NewRolePermission(ManagerLabel, true, true, true, true, 9)

	var permissions = []role.RolePermission{
		role.AnonymousPermission,
		role.RoleLessPermission,
		EmployeePermission,
		ManagerPermission,
	}

	return NewAuthorization(permissions)
}

func GetPermissionMap(permissions []role.RolePermission) map[string]role.RolePermission {
	permissionMap := make(map[string]role.RolePermission)
	for _, permission := range permissions {
		permissionMap[string(permission.Label)] = permission
	}
	return permissionMap
}

func GetPriorityRolePermission(permissions []role.RolePermission, authentic *user.Authentic) (role.RolePermission, error) {

	if authentic == nil {
		return role.AnonymousPermission, nil
	}

	if authentic.Roles == nil || len(authentic.Roles) == 0 {
		return role.RoleLessPermission, nil
	}

	permissionMap := GetPermissionMap(permisions)

	var permission role.RolePermission = nil
	for _, r := range roles {
		var roleLabel = string(r.Label)
		var p = permissionMap[roleLabel]
		if p == nil {
			return nil, utility.NewNilError("role_permission." + roleLabel, "RolePermission not found")
		}
		if permission == nil || p.Priority > permission.Priority {
			permission = p
		}
	}

	return permission, nil
}

func (auth *Authorization) Authorize(require RequirePermission, authentic *user.Authentic) error {
	var permission, err = GetPriorityRolePermission(auth.Permissions, authentic)
	if err != nil {
		return err
	}

	if require.SelfEdit && !permission.SelfEdit {
		return NewAuthorizationError(permission.RoleLabel, "SelfEdit", "Permission denied for self edit")
	}

	if require.CompanyAccess && !permission.CompanyAccess {
		return NewAuthorizationError(permission.RoleLabel, "CompanyAccess", "Permission denied for company access")
	}

	if require.CompanyInvite && !permission.CompanyInvite {
		return NewAuthorizationError(permission.RoleLabel, "CompanyInvite", "Permission denied for company invite")
	}

	if require.CompanyEdit && !permission.CompanyEdit {
		return NewAuthorizationError(permission.RoleLabel, "CompanyEdit", "Permission denied for company edit")
	}

	return nil
}
