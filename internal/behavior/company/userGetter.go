package company

import (
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	entryCompanyUser "github.com/motojouya/geezer_auth/internal/entry/transfer/companyUser"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type UserGetterDB interface {
	userQuery.GetUserAuthenticOfCompanyUserQuery
}

type UserGetter interface {
	Execute(entry entryCompanyUser.CompanyUserGetter, company shelterCompany.Company) (*shelterUser.UserAuthentic, error)
}

type UserGet struct {
	local localPkg.Localer
	db    UserGetterDB
}

func NewUserGet(local localPkg.Localer, db UserGetterDB) *UserGet {
	return &UserGet{
		local: local,
		db:    db,
	}
}

func (getter UserGet) Execute(entry entryCompanyUser.CompanyUserGetter, company shelterCompany.Company) (*shelterUser.UserAuthentic, error) {
	now := getter.local.GetNow()

	userIdentifier, err := entry.GetUserIdentifier()
	if err != nil {
		return nil, err
	}

	dbUserAuthentic, err := getter.db.GetUserAuthenticOfCompanyUser(string(company.Identifier), string(userIdentifier), now)
	if err != nil {
		return nil, err
	}

	if dbUserAuthentic == nil {
		keys := map[string]string{"company_identifier": string(company.Identifier), "user_identifier": string(userIdentifier)}
		return nil, essence.NewNotFoundError("user", keys, "user not found")
	}

	return dbUserAuthentic.ToShelterUserAuthentic()
}
