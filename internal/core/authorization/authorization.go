package authorization

import (
	"github.com/motojouya/geezer_auth/internal/model/role"
	pkgUtility "github.com/motojouya/geezer_auth/pkg/utility"
	utility "github.com/motojouya/geezer_auth/internal/utility"
	text "github.com/motojouya/geezer_auth/pkg/model/text"
	user "github.com/motojouya/geezer_auth/pkg/model/user"
	"slices"
)

type Authorization struct {
	Permisions []role.RolePermission
}

func NewAuthorization(permissions []role.RolePermission) *Authorization {
	return &Authorization{
		Permisions: permissions,
	}
}

func CreateAuthorization(permissions []role.RolePermission) Authorization {
	var perms = slices.Clone(permissions)
	perms = apend(perms, role.AnonymousPermission, role.RoleLessPermission)

	return NewAuthorization(perms)
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
		// var p, ok = utility.Find(permissions, role.PermissionIs(r.Label)) // こうも書けるが、パフォーマンス的に悪い
		if p == nil {
			return nil, pkgUtility.NewNilError("role_permission." + roleLabel, "RolePermission not found")
		}
		if permission == nil || p.Priority > permission.Priority {
			permission = p
		}
	}

	return permission, nil
}

func (auth Authorization) Authorize(require RequirePermission, authentic *user.Authentic) error {
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
