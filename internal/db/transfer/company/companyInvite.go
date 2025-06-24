package company

import (
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
	core "github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/text"
	"time"
)

type CompanyInvite struct {
	PersistKey        uint
	CompanyPersistKey uint
	Token             string
	RoleLabel         string
	RegisterDate      time.Time
	ExpireDate        time.Time
}

type CompanyInviteFull struct {
	CompanyInvite
	CompanyIdentifier     string
	CompanyName           string
	CompanyRegisteredDate time.Time
	RoleName              string
	RoleDescription       string
	RoleRegisteredDate    time.Time
}

func FromCoreCompanyInvite(invite core.CompanyInvite) CompanyInvite {
	return CompanyInvite{
		CompanyPersistKey:     invite.Company.PersistKey,
		Token:                 string(invite.Token),
		RoleLabel:             string(invite.Role.Label),
		RegisterDate:          invite.RegisterDate,
		ExpireDate:            invite.ExpireDate,
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
