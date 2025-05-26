package authorization

import (
	"github.com/motojouya/geezer_auth/internal/model/role"
	pkgUtility "github.com/motojouya/geezer_auth/pkg/utility"
	utility "github.com/motojouya/geezer_auth/internal/utility"
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

func GetPriorityRolePermission(permissions []role.RolePermission, authentic *user.Authentic) (role.RolePermission, error) {

	if authentic == nil {
		return role.AnonymousPermission, nil
	}

	if authentic.Roles == nil || len(authentic.Roles) == 0 {
		return role.RoleLessPermission, nil
	}

	var permissionMap = utility.ToMap(permisions, role.PermissionKey)

	var permission role.RolePermission = nil
	for _, r := range roles {
		var roleLabel = string(r.Label)
		var p = permissionMap[roleLabel]
		// var p = utility.Find(permissions, role.PermissionWhen(r.Label)) // こうも書けるが、パフォーマンス的に悪い
		if p == nil {
			return nil, pkgUtility.NewNilError("role_permission." + roleLabel, "RolePermission not found")
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
