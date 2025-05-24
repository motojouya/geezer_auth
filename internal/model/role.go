package model

import (
	"time"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
)

/*
 * Roleは管理者が登録する想定なので、基本的には削除されない
 */
type UnsavedRole struct {
	pkg.Role
	Description    Text
	RegisteredDate time.Time
}

type Role struct {
	RoleId uint
	UnsavedRole
}

func CreateRole(name pkg.Name, label pkg.Label, description Text, registeredDate time.Time) UnsavedRole {
	return UnsavedRole{
		Role:           pkg.NewRole(label, name),
		Description:    description,
		RegisteredDate: registeredDate,
	}
}

func NewRole(roleId uint, name pkg.Name, label pkg.Label, description Text, registeredDate time.Time) Role {
	return Role{
		RoleId:      roleId,
		UnsavedRole: UnsavedRole{
			Role:           pkg.NewRole(label, name),
			Description:    description,
			RegisteredDate: registeredDate,
		},
	}
}
