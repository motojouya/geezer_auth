package model

import (
	"time"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
)

// SelfEditがないと、自身の情報を編集するのがだれでもできるとか、権限で管理ができなくなるので
// 逆にだれでもできるユーザ登録とかは、特に権限を設定する必要がない。
type RolePermission struct {
	RoleLabel     pkg.Label
	SelfEdit      bool
	CompanyAccess bool
	CompanyInvite bool
	CompanyEdit   bool
	Priority      uint
}

func NewRolePermission(roleLabel pkg.Label, selfEdit bool, companyAccess bool, companyInvite bool, companyEdit bool, priority uint) RolePermission {
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
const RoleLessLabel = pkg.NewLabel("ROLE_LESS")
const RoleLessPermission = NewRolePermission(RoleLessLabel, true, false, false, false, 1)

// そもそも認証がない利用者ができることを定義する
const AnonymousLabel = pkg.NewLabel("ANONYMOUS")
const AnonymousPermission = NewRolePermission(AnonymousLabel, false, false, false, false, 0)
