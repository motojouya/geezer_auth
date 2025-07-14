package role

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	core "github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	text "github.com/motojouya/geezer_auth/pkg/core/text"
)

type RolePermission struct {
	RoleLabel     string `db:"role_label,primarykey"`
	SelfEdit      bool   `db:"self_edit"`
	CompanyAccess bool   `db:"company_access"`
	CompanyInvite bool   `db:"company_invite"`
	CompanyEdit   bool   `db:"company_edit"`
	Priority      uint   `db:"priority"`
}

func AddRolePermissionTable(dbMap *gorp.DbMap) {
	dbMap.AddTable(RolePermission{})
}

var SelectRolePermission = utility.Dialect.From("role_permission").As("rp").Select(
	goqu.C("rp.role_label").As("role_label"),
	goqu.C("rp.self_edit").As("self_edit"),
	goqu.C("rp.company_access").As("company_access"),
	goqu.C("rp.company_invite").As("company_invite"),
	goqu.C("rp.company_edit").As("company_edit"),
	goqu.C("rp.priority").As("priority"),
)

func FromCoreRolePermission(r core.RolePermission) RolePermission {
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

// testdata投入時に楽するためのもの。アプリケーションからは利用を想定しない。
func NewRolePermission(roleLabel string, selfEdit bool, companyAccess bool, companyInvite bool, companyEdit bool, priority uint) RolePermission {
	return RolePermission{
		RoleLabel:     roleLabel,
		SelfEdit:      selfEdit,
		CompanyAccess: companyAccess,
		CompanyInvite: companyInvite,
		CompanyEdit:   companyEdit,
		Priority:      priority,
	}
}
