package company

import (
	authBehavior "github.com/motojouya/geezer_auth/internal/behavior/authorization"
	companyBehavior "github.com/motojouya/geezer_auth/internal/behavior/company"
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	entryCompany "github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/authorization"
	essence "github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
)

type GetCompanyControl struct {
	essence.Closable
	authorization *authorization.Authorization
	companyGetter companyBehavior.CompanyGetter
}

func NewGetCompanyControl(
	database essence.Closable,
	authorization *authorization.Authorization,
	companyGetter companyBehavior.CompanyGetter,
) *GetCompanyControl {
	return &GetCompanyControl{
		Closable:      database,
		authorization: authorization,
		companyGetter: companyGetter,
	}
}

func CreateGetCompanyControl() (*GetCompanyControl, error) {
	var env = localPkg.CreateEnvironment()

	database, err := configBehavior.NewDatabaseGet(env).GetDatabase()
	if err != nil {
		return nil, err
	}

	authorization, err := authBehavior.NewAuthorizationGet(database).GetAuthorization()
	if err != nil {
		return nil, err
	}

	companyGetter := companyBehavior.NewCompanyGet(database)

	return NewGetCompanyControl(database, authorization, companyGetter), nil
}

var getCompanyPermission = shelterRole.NewRequirePermission(true, false, false, false)

func GetCompanyExecute(control *GetCompanyControl, entry entryCompany.CompanyGetRequest, authentic *pkgUser.Authentic) (entryCompany.CompanyGetResponse, error) {

	if err := control.authorization.Authorize(getCompanyPermission, authentic); err != nil {
		return entryCompany.CompanyGetResponse{}, err
	}

	company, err := control.companyGetter.Execute(entry)
	if err != nil {
		return entryCompany.CompanyGetResponse{}, err
	}
	if company == nil {
		keys := map[string]string{"identifier": entry.CompanyGet.Identifier}
		return entryCompany.CompanyGetResponse{}, essence.NewNotFoundError("company", keys, "company not found")
	}

	return entryCompany.FromShelterCompany(*company), nil
}
