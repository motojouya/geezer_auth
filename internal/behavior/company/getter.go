package company

import (
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	companyQuery "github.com/motojouya/geezer_auth/internal/db/query/company"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type CompanyGetterDB interface {
	companyQuery.GetCompanyQuery
}

type CompanyGetter interface {
	Execute(entry entryCompany.CompanyGetter) (*shelterCompany.Company, error)
}

type CompanyGet struct {
	db    CompanyGetterDB
}

func NewCompanyGet(db CompanyGetterDB) CompanyGet {
	return &CompanyGet{
		db:    db,
	}
}

func (getter CompanyGet) Execute(entry entryCompany.CompanyGetter) (*shelterCompany.Company, error) {

	identifier, err := entry.GetCompanyIdentifier()
	if err != nil {
		return nil, err
	}

	dbCompany, err := getter.db.GetCompany(string(identifier))
	if err != nil {
		return nil, err
	}

	if dbCompany == nil {
		keys := map[string]string{"identifier": string(identifier)}
		return nil, essence.NewNotFoundError("company", keys, "company not found")
	}

	return dbCompany.ToShelterCompany()
}
