package company

import (
	"github.com/go-gorp/gorp"
	dbCompany "github.com/motojouya/geezer_auth/internal/db/transfer/company"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
)

type InviteTokenIssuerDB interface {
	gorp.SqlExecutor
}

type InviteTokenIssuer interface {
	Execute(company shelterCompany.Company, role shelterRole.Role) (shelterText.Token, error)
}

type InviteTokenIssue struct {
	local localPkg.Localer
	db    InviteTokenIssuerDB
}

func NewInviteTokenIssue(local localPkg.Localer, database InviteTokenIssuerDB) *InviteTokenIssue {
	return &InviteTokenIssue{
		db:    database,
		local: local,
	}
}

func (issuer InviteTokenIssue) Execute(company shelterCompany.Company, role shelterRole.Role) (shelterText.Token, error) {
	now := issuer.local.GetNow()

	inviteTokenSource, err := issuer.local.GenerateUUID()
	if err != nil {
		return shelterText.Token(""), err
	}

	inviteToken, err := shelterText.CreateToken(inviteTokenSource)
	if err != nil {
		return shelterText.Token(""), err
	}

	companyInvite := shelterCompany.CreateCompanyInvite(company, inviteToken, role, now)
	dbCompanyInvite := dbCompany.FromShelterCompanyInvite(companyInvite)

	if err := issuer.db.Insert(&dbCompanyInvite); err != nil {
		return shelterText.Token(""), err
	}

	return inviteToken, nil
}
