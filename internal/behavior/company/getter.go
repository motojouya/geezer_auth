package company

import (
	companyQuery "github.com/motojouya/geezer_auth/internal/db/query/company"
	entryCompany "github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
)

type CompanyGetterDB interface {
	companyQuery.GetCompanyQuery
}

type CompanyGetter interface {
	Execute(entry entryCompany.CompanyGetter) (shelterCompany.Company, error)
}

type CompanyGet struct {
	db CompanyGetterDB
}

func NewCompanyGet(db CompanyGetterDB) *CompanyGet {
	return &CompanyGet{
		db: db,
	}
}

func (getter CompanyGet) Execute(entry entryCompany.CompanyGetter) (shelterCompany.Company, error) {

	identifier, err := entry.GetCompanyIdentifier()
	if err != nil {
		return shelterCompany.Company{}, err
	}

	dbCompany, err := getter.db.GetCompany(string(identifier))
	if err != nil {
		return shelterCompany.Company{}, err
	}

	if dbCompany == nil {
		keys := map[string]string{"identifier": string(identifier)}
		return shelterCompany.Company{}, essence.NewNotFoundError("company", keys, "company not found")
	}

	return dbCompany.ToShelterCompany()
}
