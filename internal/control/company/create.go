package company

import (
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	companyBehavior "github.com/motojouya/geezer_auth/internal/behavior/company"
	roleBehavior "github.com/motojouya/geezer_auth/internal/behavior/role"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/motojouya/geezer_auth/internal/control/utility"
	"github.com/motojouya/geezer_auth/internal/db"
	entryCompany "github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
)

type CreateControl struct {
	db.TransactionalDatabase
	Authorization  *authorization.Authorization
	CompanyCreator companyBehavior.UserCreator
	RoleGetter     roleBehavior.RoleGetter
	RoleAssigner   companyBehavior.RoleAssigner
}

func NewCreateControl(
	database db.TransactionalDatabase,
	authorization *authorization.Authorization,
	companyCreator companyBehavior.UserCreator,
	roleGetter roleBehavior.RoleGetter,
	roleAssigner companyBehavior.RoleAssigner,
) *CreateControl {
	return &CreateControl{
		TransactionalDatabase: database,
		Authorization:	       authorization,
		CompanyCreator:        companyBehavior.UserCreator,
		RoleGetter:            roleBehavior.RoleGetter,
		RoleAssigner:          companyBehavior.RoleAssigner,
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

	companyCreator := companyBehavior.NewCompanyCreate(local, database)
	roleGetter := roleBehavior.NewRoleGet(database)
	roleAssigner := companyBehavior.NewRoleAssign(local, database)

	return NewCreateControl(database, authorization, companyCreator, roleGetter, roleAssigner), nil
}

var createCompanyPermission = shelterRole.NewRequirePermission(true, false, false, false)

var CreateExecute = utility.Transact(func(control *CreateControl, entry entryCompany.CompanyCreateRequest, authentic *pkgUser.Authentic) (*entryCompany.CompanyGetResponse, error) {

	if err := control.authorization.Authorize(createCompanyPermission, authentic); err != nil {
		return nil, err
	}

	company, err := control.companyCreator.Execute(entry)
	if err != nil {
		return nil, err
	}

	role, err := control.roleGetter.Execute(struct {
		entryCompany.RoleGetter
	}{
		GetRoleLabel: func() (pkgText.Label, error) {
			return shelterRole.RoleAdminLabel, nil
		},
	})
	if err != nil {
		return nil, err
	}

	if _, err = control.roleAssigner.Execute(company, role); err != nil {
		return nil, err
	}

	return entryCompany.FromShelterCompany(company), nil
})
