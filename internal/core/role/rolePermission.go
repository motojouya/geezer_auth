package role

import (
	text "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

// SelfEditがないと、自身の情報を編集するのがだれでもできるとか、権限で管理ができなくなるので
// 逆にだれでもできるユーザ登録とかは、特に権限を設定する必要がない。
type RolePermission struct {
	RoleLabel     text.Label
	SelfEdit      bool
	CompanyAccess bool
	CompanyInvite bool
	CompanyEdit   bool
	Priority      uint
}

func NewRolePermission(roleLabel text.Label, selfEdit bool, companyAccess bool, companyInvite bool, companyEdit bool, priority uint) RolePermission {
	return RolePermission{
		RoleLabel:     roleLabel,
		SelfEdit:      selfEdit,
		CompanyAccess: companyAccess,
		CompanyInvite: companyInvite,
		CompanyEdit:   companyEdit,
		Priority:      priority,
	}
}

// RoleはCompanyにAssignされて、そのCompanyで適用されるものなので、ロールのない状態でできる権限を設定する
const RoleLessLabel = text.NewLabel("ROLE_LESS")
const RoleLessPermission = NewRolePermission(RoleLessLabel, true, false, false, false, 1)

// そもそも認証がない利用者ができることを定義する
const AnonymousLabel = text.NewLabel("ANONYMOUS")
const AnonymousPermission = NewRolePermission(AnonymousLabel, false, false, false, false, 0)

func PermissionKey(permission RolePermission) string {
	return string(permission.RoleLabel)
}

func PermissionIs(label text.Label) func(permission RolePermission) bool {
	return func(permission RolePermission) bool {
		return string(permission.RoleLabel) == string(label)
	}
}
