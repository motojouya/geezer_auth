package model

import (
	"time"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
)

// SelfEditがないと、自身の情報を編集するのがだれでもできるとか、権限で管理ができなくなるので
// 逆にだれでもできるユーザ登録とかは、特に権限を設定する必要がない。
type RolePermission struct {
	RoleLabel      pkg.Label
	SelfEdit       bool
	CompanyAccess  bool
	CompanyInvite  bool
	CompanyEdit    bool
}

func NewRolePermission(roleLabel pkg.Label, selfEdit bool, companyAccess bool, companyInvite bool, companyEdit bool) RolePermission {
	return RolePermission{
		RoleLabel:     roleLabel,
		SelfEdit:      selfEdit,
		CompanyAccess: companyAccess,
		CompanyInvite: companyInvite,
		CompanyEdit:   companyEdit,
	}
}

const DefaultRoleLabel = pkg.NewLabel("DEFAULT")
const DefaultRolePermission = NewRolePermission(DefaultRoleLabel, true, false, false, false)
