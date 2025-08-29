package company

import (
	companyQuery "github.com/motojouya/geezer_auth/internal/db/query/company"
	entryCompany "github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
)

type CompanyGetterDB interface {
	companyQuery.GetCompanyQuery
}

type CompanyGetter interface {
	Execute(entry entryCompany.CompanyGetter) (*shelterCompany.Company, error)
}

type CompanyGet struct {
	db CompanyGetterDB
}

func NewCompanyGet(db CompanyGetterDB) *CompanyGet {
	return &CompanyGet{
		db: db,
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
		return nil, nil

	}

	resultCompany, err := dbCompany.ToShelterCompany()
	if err != nil {
		return nil, err
	}

	return &resultCompany, nil
}
