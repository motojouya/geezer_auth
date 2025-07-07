package authorization

import (
	"github.com/motojouya/geezer_auth/internal/core/essence"
	"github.com/motojouya/geezer_auth/internal/core/role"
	pkgEssence "github.com/motojouya/geezer_auth/pkg/core/essence"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"slices"
)

type Authorization struct {
	Permissions []role.RolePermission
}

func NewAuthorization(permissions []role.RolePermission) *Authorization {
	return &Authorization{
		Permissions: permissions,
	}
}

func CreateAuthorization(permissions []role.RolePermission) *Authorization {
	var perms = slices.Clone(permissions)
	perms = append(perms, role.AnonymousPermission, role.RoleLessPermission)

	return NewAuthorization(perms)
}

func GetPriorityRolePermission(permissions []role.RolePermission, authentic *user.Authentic) (role.RolePermission, error) {

	if authentic == nil {
		return role.AnonymousPermission, nil
	}

	if authentic.User.CompanyRole == nil || len(authentic.User.CompanyRole.Roles) == 0 {
		return role.RoleLessPermission, nil
	}

	var permissionMap = essence.ToMap(permissions, role.PermissionKey)

	var permission *role.RolePermission = nil
	for _, r := range authentic.User.CompanyRole.Roles {
		var roleLabel = string(r.Label)
		var p, exists = permissionMap[roleLabel]
		// var p, ok = essence.Find(permissions, role.PermissionIs(r.Label)) // こうも書けるが、パフォーマンス的に悪い
		if !exists {
			return role.RolePermission{}, pkgEssence.NewNilError("role_permission."+roleLabel, "RolePermission not found")
		}
		if permission == nil || p.Priority > permission.Priority {
			permission = &p
		}
	}

	return *permission, nil
}

func (auth Authorization) Authorize(require role.RequirePermission, authentic *user.Authentic) error {
	var permission, err = GetPriorityRolePermission(auth.Permissions, authentic)
	if err != nil {
		return err
	}

	if require.SelfEdit && !permission.SelfEdit {
		return NewAuthorizationError(string(permission.RoleLabel), "SelfEdit", "Permission denied for self edit")
	}

	if require.CompanyAccess && !permission.CompanyAccess {
		return NewAuthorizationError(string(permission.RoleLabel), "CompanyAccess", "Permission denied for company access")
	}

	if require.CompanyInvite && !permission.CompanyInvite {
		return NewAuthorizationError(string(permission.RoleLabel), "CompanyInvite", "Permission denied for company invite")
	}

	if require.CompanyEdit && !permission.CompanyEdit {
		return NewAuthorizationError(string(permission.RoleLabel), "CompanyEdit", "Permission denied for company edit")
	}

	return nil
}
