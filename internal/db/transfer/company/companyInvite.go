package company

import (
	"github.com/doug-martin/goqu/v9"
	core "github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"time"
)

type CompanyInvite struct {
	PersistKey        uint      `db:"persist_key,primarykey,autoincrement"`
	CompanyPersistKey uint      `db:"company_persist_key"`
	Token             string    `db:"verify_token"`
	RoleLabel         string    `db:"role_label"`
	RegisterDate      time.Time `db:"register_date"`
	ExpireDate        time.Time `db:"expire_date"`
}

type CompanyInviteFull struct {
	CompanyInvite
	CompanyIdentifier     string    `db:"company_identifier"`
	CompanyName           string    `db:"company_name"`
	CompanyRegisteredDate time.Time `db:"company_register_date"`
	RoleName              string    `db:"role_name"`
	RoleDescription       string    `db:"role_description"`
	RoleRegisteredDate    time.Time `db:"role_register_date"`
}

// var SelectCompanyInvite = utility.Dialect.From("company_invite").As("ci").Select(
// 	goqu.C("ci.persist_key").As("persist_key"),
// 	goqu.C("ci.company_persist_key").As("company_persist_key"),
// 	goqu.C("ci.verify_token").As("verify_token"),
// 	goqu.C("ci.role_label").As("role_label"),
// 	goqu.C("ci.register_date").As("register_date"),
// 	goqu.C("ci.expire_date").As("expire_date"),
// )

var SelectCompanyInvite = utility.Dialect.From("company_invite").As("ci").InnerJoin(
	goqu.T("company").As("c"),
	goqu.On(goqu.Ex{"ci.company_persist_key": goqu.I("c.persist_key")}),
).InnerJoin(
	goqu.T("role").As("r"),
	goqu.On(goqu.Ex{"ci.role_label": goqu.I("r.label")}),
).Select(
	goqu.C("ci.persist_key").As("persist_key"),
	goqu.C("ci.company_persist_key").As("company_persist_key"),
	goqu.C("c.identifier").As("company_identifier"),
	goqu.C("c.name").As("company_name"),
	goqu.C("c.register_date").As("company_register_date"),
	goqu.C("ci.verify_token").As("verify_token"),
	goqu.C("ci.role_label").As("role_label"),
	goqu.C("r.name").As("role_name"),
	goqu.C("r.description").As("role_description"),
	goqu.C("r.register_date").As("role_register_date"),
	goqu.C("ci.register_date").As("register_date"),
	goqu.C("ci.expire_date").As("expire_date"),
)

func FromCoreCompanyInvite(invite core.CompanyInvite) CompanyInvite {
	return CompanyInvite{
		CompanyPersistKey: invite.Company.PersistKey,
		Token:             string(invite.Token),
		RoleLabel:         string(invite.Role.Label),
		RegisterDate:      invite.RegisterDate,
		ExpireDate:        invite.ExpireDate,
	}
}

func (c CompanyInviteFull) ToCoreCompanyInvite() (core.CompanyInvite, error) {
	var company, companyErr = (Company{
		PersistKey:     c.CompanyPersistKey,
		Identifier:     c.CompanyIdentifier,
		Name:           c.CompanyName,
		RegisteredDate: c.CompanyRegisteredDate,
	}).ToCoreCompany()
	if companyErr != nil {
		return core.CompanyInvite{}, companyErr
	}

	var role, roleErr = (role.Role{
		Label:          c.RoleLabel,
		Name:           c.RoleName,
		Description:    c.RoleDescription,
		RegisteredDate: c.RoleRegisteredDate,
	}).ToCoreRole()
	if roleErr != nil {
		return core.CompanyInvite{}, roleErr
	}

	var token, tokenErr = text.NewToken(c.Token)
	if tokenErr != nil {
		return core.CompanyInvite{}, tokenErr
	}

	return core.NewCompanyInvite(
		c.PersistKey,
		company,
		token,
		role,
		c.RegisterDate,
		c.ExpireDate,
	), nil
}
