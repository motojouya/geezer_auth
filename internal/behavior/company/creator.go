package company

import (
	"github.com/go-gorp/gorp"
	companyQuery "github.com/motojouya/geezer_auth/internal/db/query/company"
	dbCompany "github.com/motojouya/geezer_auth/internal/db/transfer/company"
	entryCompany "github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type CompanyCreatorDB interface {
	gorp.SqlExecutor
	companyQuery.GetCompanyQuery
}

type CompanyCreator interface {
	Execute(entry entryCompany.CompanyCreator) (shelterCompany.Company, error)
}

type CompanyCreate struct {
	local localPkg.Localer
	db    CompanyCreatorDB
}

func NewCompanyCreate(local localPkg.Localer, db CompanyCreatorDB) *CompanyCreate {
	return &CompanyCreate{
		local: local,
		db:    db,
	}
}

func createCompanyIdentifier(local localPkg.Localer) func() (pkgText.Identifier, error) {
	return func() (pkgText.Identifier, error) {
		var ramdomString = local.GenerateRamdomString(pkgText.IdentifierLength, pkgText.IdentifierChar)
		var identifier, err = shelterCompany.CreateCompanyIdentifier(ramdomString)
		if err != nil {
			return pkgText.Identifier(""), err
		}
		return identifier, nil
	}
}

func checkCompanyIdentifier(companyCreatorDB CompanyCreatorDB) func(pkgText.Identifier) (bool, error) {
	return func(identifier pkgText.Identifier) (bool, error) {
		var company, err = companyCreatorDB.GetCompany(string(identifier))
		if err != nil {
			return false, err
		}
		return company == nil, nil
	}
}

func (creator CompanyCreate) Execute(entry entryCompany.CompanyCreator) (shelterCompany.Company, error) {
	now := creator.local.GetNow()

	companyIdentifier, err := shelterText.GetString(createCompanyIdentifier(creator.local), checkCompanyIdentifier(creator.db), 10)
	if err != nil {
		return shelterCompany.Company{}, err
	}

	unsavedCompany, err := entry.ToShelterCompany(companyIdentifier, now)
	if err != nil {
		return shelterCompany.Company{}, err
	}

	var dbCompanyValue = dbCompany.FromShelterCompany(unsavedCompany)

	if err = creator.db.Insert(&dbCompanyValue); err != nil {
		return shelterCompany.Company{}, err
	}

	dbCompanyResult, err := creator.db.GetCompany(string(companyIdentifier))
	if err != nil {
		return shelterCompany.Company{}, err
	}

	if dbCompanyResult == nil {
		keys := map[string]string{"identifier": string(companyIdentifier)}
		return shelterCompany.Company{}, essence.NewNotFoundError("company", keys, "company not found")
	}

	return dbCompanyResult.ToShelterCompany()
}
