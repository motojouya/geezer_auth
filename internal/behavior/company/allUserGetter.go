package company

import (
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type AllUserGetterDB interface {
	userQuery.GetUserAuthenticOfCompanyQuery
}

type AllUserGetter interface {
	Execute(entry entryCompany.CompanyGetter) ([]shelterUser.UserAuthentic, error)
}

type AllUserGet struct {
	local localPkg.Localer
	db    AllUserGetterDB
}

func NewAllUserGet(local localPkg.Localer, db UserGetterDB) *AllUserGet {
	return &AllUserGet{
		local: local,
		db:    db,
	}
}

func (getter AllUserGet) Execute(entry entryCompany.CompanyGetter) ([]shelterUser.UserAuthentic, error) {
	now := getter.local.GetNow()

	companyIdentifier, err := entry.GetCompanyIdentifier()
	if err != nil {
		return nil, err
	}

	dbUserAuthentics, err := getter.db.GetUserAuthenticOfCompany(string(companyIdentifier), now)
	if err != nil {
		return nil, err
	}

	if len(dbUserAuthentics) == 0 {
		keys := map[string]string{"company_identifier": string(companyIdentifier)}
		return nil, essence.NewNotFoundError("user", keys, "user not found")
	}

	var userAuthentics []shelterUser.UserAuthentic
	for _, dbUserAuthentic := range dbUserAuthentics {
		userAuthentic, err := dbUserAuthentic.ToShelterUserAuthentic()
		if err != nil {
			return nil, err
		}
		userAuthentics = append(userAuthentics, userAuthentic)
	}

	return userAuthentics, nil
}
