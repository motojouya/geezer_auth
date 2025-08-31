package company

import (
	authBehavior "github.com/motojouya/geezer_auth/internal/behavior/authorization"
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	companyBehavior "github.com/motojouya/geezer_auth/internal/behavior/company"
	entryCommon "github.com/motojouya/geezer_auth/internal/entry/transfer/common"
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
	var local = localPkg.CreateLocal()
	var env = localPkg.CreateEnvironment()

	database, err := configBehavior.NewDatabaseGet(env).GetDatabase()
	if err != nil {
		return nil, err
	}

	authorization, err := authBehavior.NewAuthorizationGet(database).GetAuthorization()
	if err != nil {
		return nil, err
	}

	companyGetter := companyBehavior.NewCompanyGet(local, database)

	return NewGetCompanyControl(database, authorization, companyGetter), nil
}

var getCompanyPermission = shelterRole.NewRequirePermission(true, false, false, false)

func GetCompanyExecute(control *GetCompanyControl, entry entryCompany.CompanyGetRequest, authentic *pkgUser.Authentic) (*entryCompany.CompanyGetResponse, error) {

	if err := control.authorization.Authorize(getCompanyPermission, authentic); err != nil {
		return nil, err
	}

	company, err := control.companyGetter.Execute(entry)
	if err != nil {
		return nil, err
	}

	return entryCompany.FromShelterCompany(company), nil
}
