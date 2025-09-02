package company

import (
	authBehavior "github.com/motojouya/geezer_auth/internal/behavior/authorization"
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	companyBehavior "github.com/motojouya/geezer_auth/internal/behavior/company"
	entryCompany "github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/authorization"
	essence "github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
)

type GetUserControl struct {
	essence.Closable
	authorization *authorization.Authorization
	allUserGetter companyBehavior.AllUserGetter
}

func NewGetUserControl(
	database essence.Closable,
	authorization *authorization.Authorization,
	allUserGetter companyBehavior.AllUserGetter,
) *GetUserControl {
	return &GetUserControl{
		Closable:      database,
		authorization: authorization,
		allUserGetter: allUserGetter,
	}
}

func CreateGetUserControl() (*GetUserControl, error) {
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

	allUserGetter := companyBehavior.NewAllUserGet(local, database)

	return NewGetUserControl(database, authorization, allUserGetter), nil
}

var getUserPermission = shelterRole.NewRequirePermission(true, true, false, false)

func GetUserExecute(control *GetUserControl, entry entryCompany.CompanyGetRequest, authentic *pkgUser.Authentic) (*entryCompany.CompanyUserResponse, error) {

	if err := control.authorization.Authorize(getUserPermission, authentic); err != nil {
		return nil, err
	}

	users, err := control.allUserGetter.Execute(entry)
	if err != nil {
		return nil, err
	}

	return entryCompany.FromShelterUserAuthentic(users), nil
}
