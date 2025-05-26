package authorization

import (
	"github.com/motojouya/geezer_auth/internal/model/role"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
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

// TODO DBアクセスしてロードする
func CreateAuthorization() *Authorization {
	var EmployeeLabel = pkg.NewLabel("EMPLOYEE")
	var EmployeePermission = role.NewRolePermission(EmployeeLabel, true, true, false, false, 5)

	var ManagerLabel = pkg.NewLabel("MANAGER")
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

func GetPriorityRolePermission(permissions []role.RolePermission, authentic *pkg.Authentic) (role.RolePermission, error) {

	if authentic == nil {
		return role.AnonymousPermission, nil
	}

	if authentic.Roles == nil || len(authentic.Roles) == 0 {
		return role.RoleLessPermission, nil
	}

	permissionMap := GetPermissionMap(permisions)

	var permission role.RolePermission = nil
	for _, r := range roles {
		var p = permissionMap[string(r.Label)]
		if p == nil {
			// TODO System Config Errorみたいな感じのを定義したい
			return nil, error.Errorf("RolePermission not found: %s", r.Label)
		}
		if permission == nil || p.Priority > permission.Priority {
			permission = p
		}
	}

	return permission, nil
}

// 分かりづらいが、認可エラーと、設定エラーが出るので、呼び出し側で判断すること
func (auth *Authorization) Authorize(require RequirePermission, authentic *pkg.Authentic) error {
	var permission, err = GetPriorityRolePermission(auth.Permissions, authentic)
	if err != nil {
		// ここでのエラーは設定エラーのため、呼び出し側でのハンドリングが違うはず
		return err
	}

	if require.SelfEdit && !permission.SelfEdit {
		// TODO error type
		return ErrPermissionDenied
	}

	if require.CompanyAccess && !permission.CompanyAccess {
		// TODO error type
		return ErrPermissionDenied
	}

	if require.CompanyInvite && !permission.CompanyInvite {
		// TODO error type
		return ErrPermissionDenied
	}

	if require.CompanyEdit && !permission.CompanyEdit {
		// TODO error type
		return ErrPermissionDenied
	}

	return nil
}
