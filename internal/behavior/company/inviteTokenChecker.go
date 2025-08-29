package company

import (
	companyQuery "github.com/motojouya/geezer_auth/internal/db/query/company"
	entryCompanyUser "github.com/motojouya/geezer_auth/internal/entry/transfer/companyUser"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	shelterAuth "github.com/motojouya/geezer_auth/internal/shelter/authorization"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	essence "github.com/motojouya/geezer_auth/internal/shelter/essence"
)

type InviteTokenCheckerDB interface {
	companyQuery.GetCompanyInviteQuery
}

type InviteTokenChecker interface {
	Execute(entry entryCompanyUser.InviteTokenGetter, company shelterCompany.Company) (shelterRole.Role, error)
}

type InviteTokenCheck struct {
	local localPkg.Localer
	db    InviteTokenCheckerDB
}

func NewInviteTokenCheck(local localPkg.Localer, database InviteTokenCheckerDB) *InviteTokenCheck {
	return &InviteTokenCheck{
		db:    database,
		local: local,
	}
}

func (checker InviteTokenCheck) Execute(entry entryCompanyUser.InviteTokenGetter, company shelterCompany.Company) (shelterRole.Role, error) {
	now := checker.local.GetNow()

	inviteToken, err := entry.GetToken()
	if err != nil {
		return shelterRole.Role{}, err
	}

	dbCompanyInvite, err := checker.db.GetCompanyInvite(string(company.Identifier), string(inviteToken))
	if err != nil {
		return shelterRole.Role{}, err
	}

	if dbCompanyInvite == nil {
		keys := map[string]string{"token": string(inviteToken), "identifier": string(company.Identifier)}
		return shelterRole.Role{}, essence.NewNotFoundError("company_invite", keys, "company_invite not found")
	}

	companyInvite, err := dbCompanyInvite.ToShelterCompanyInvite()

	if companyInvite.ExpireDate.Before(now) {
		return shelterRole.Role{}, shelterAuth.NewTokenExpiredError(companyInvite.ExpireDate, "invite token is expired")
	}

	return companyInvite.Role, nil
}
