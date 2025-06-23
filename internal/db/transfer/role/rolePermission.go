package role

import (
	text "github.com/motojouya/geezer_auth/pkg/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/role"
)

type RolePermission struct {
	RoleLabel     string
	SelfEdit      bool
	CompanyAccess bool
	CompanyInvite bool
	CompanyEdit   bool
	Priority      uint
}

func FromCoreRolePermission(r core.RolePermission) (RolePermission) {
	return RolePermission{
		RoleLabel:     string(r.RoleLabel),
		SelfEdit:      r.SelfEdit,
		CompanyAccess: r.CompanyAccess,
		CompanyInvite: r.CompanyInvite,
		CompanyEdit:   r.CompanyEdit,
		Priority:      r.Priority,
	}
}

func (r RolePermission) ToCoreRolePermission() (core.RolePermission, error) {
	var label, err = text.NewLabel(r.RoleLabel)
	if err != nil {
		return core.RolePermission{}, err
	}

	return core.NewRolePermission(
		label,
		r.SelfEdit,
		r.CompanyAccess,
		r.CompanyInvite,
		r.CompanyEdit,
		r.Priority,
	), nil
}
