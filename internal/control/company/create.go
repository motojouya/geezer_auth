package company

import (
	authBehavior "github.com/motojouya/geezer_auth/internal/behavior/authorization"
	companyBehavior "github.com/motojouya/geezer_auth/internal/behavior/company"
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	roleBehavior "github.com/motojouya/geezer_auth/internal/behavior/role"
	userBehavior "github.com/motojouya/geezer_auth/internal/behavior/user"
	"github.com/motojouya/geezer_auth/internal/control/utility"
	"github.com/motojouya/geezer_auth/internal/db"
	entryCompany "github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/authorization"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
)

type CreateControl struct {
	db.TransactionalDatabase
	authorization     *authorization.Authorization
	companyCreator    companyBehavior.CompanyCreator
	roleGetter        roleBehavior.RoleGetter
	userGetter        userBehavior.UserGetter
	roleAssigner      companyBehavior.RoleAssigner
	accessTokenIssuer userBehavior.AccessTokenIssuer
}

func NewCreateControl(
	database db.TransactionalDatabase,
	authorization *authorization.Authorization,
	companyCreator companyBehavior.CompanyCreator,
	roleGetter roleBehavior.RoleGetter,
	userGetter userBehavior.UserGetter,
	roleAssigner companyBehavior.RoleAssigner,
	accessTokenIssuer userBehavior.AccessTokenIssuer,
) *CreateControl {
	return &CreateControl{
		TransactionalDatabase: database,
		authorization:         authorization,
		companyCreator:        companyCreator,
		roleGetter:            roleGetter,
		userGetter:            userGetter,
		roleAssigner:          roleAssigner,
		accessTokenIssuer:     accessTokenIssuer,
	}
}

func CreateCreateControl() (*CreateControl, error) {
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

	jwtHandler, err := configBehavior.NewJwtHandlerGet(env).GetJwtHandler()
	if err != nil {
		return nil, err
	}

	companyCreator := companyBehavior.NewCompanyCreate(local, database)
	roleGetter := roleBehavior.NewRoleGet(database)
	userGetter := userBehavior.NewUserGet(local, database)
	roleAssigner := companyBehavior.NewRoleAssign(local, database)
	accessTokenIssuer := userBehavior.NewAccessTokenIssue(local, database, jwtHandler)

	return NewCreateControl(database, authorization, companyCreator, roleGetter, userGetter, roleAssigner, accessTokenIssuer), nil
}

var createCompanyPermission = shelterRole.NewRequirePermission(true, false, false, false)

type RoleGetEntry struct {
	label pkgText.Label
}

func (entry *RoleGetEntry) GetRoleLabel() (pkgText.Label, error) {
	return entry.label, nil
}

var CreateExecute = utility.Transact(func(control *CreateControl, entry entryCompany.CompanyCreateRequest, authentic *pkgUser.Authentic) (entryCompany.CompanyTokenResponse, error) {

	if err := control.authorization.Authorize(createCompanyPermission, authentic); err != nil {
		return entryCompany.CompanyTokenResponse{}, err
	}

	company, err := control.companyCreator.Execute(entry)
	if err != nil {
		return entryCompany.CompanyTokenResponse{}, err
	}

	role, err := control.roleGetter.Execute(&RoleGetEntry{
		label: shelterRole.RoleAdminLabel,
	})
	if err != nil {
		return entryCompany.CompanyTokenResponse{}, err
	}
	if role == nil {
		keys := map[string]string{"label": string(shelterRole.RoleAdminLabel)}
		return entryCompany.CompanyTokenResponse{}, essence.NewNotFoundError("role", keys, "role not found")
	}

	userAuthentic, err := control.userGetter.Execute(authentic.User.Identifier)
	if err != nil {
		return entryCompany.CompanyTokenResponse{}, err
	}

	if _, err = control.roleAssigner.Execute(company, userAuthentic, *role); err != nil {
		return entryCompany.CompanyTokenResponse{}, err
	}

	accessToken, err := control.accessTokenIssuer.Execute(userAuthentic)
	if err != nil {
		return entryCompany.CompanyTokenResponse{}, err
	}

	return entryCompany.FromShelterCompanyToken(company, accessToken), nil
})
