package companyUser

import (
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	authBehavior "github.com/motojouya/geezer_auth/internal/behavior/authorization"
	roleBehavior "github.com/motojouya/geezer_auth/internal/behavior/role"
	companyBehavior "github.com/motojouya/geezer_auth/internal/behavior/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	"github.com/motojouya/geezer_auth/internal/control/utility"
	"github.com/motojouya/geezer_auth/internal/db"
	entryCompanyUser "github.com/motojouya/geezer_auth/internal/entry/transfer/companyUser"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
	"github.com/motojouya/geezer_auth/internal/shelter/authorization"
)

type AssignControl struct {
	db.TransactionalDatabase
	authorization      *authorization.Authorization
	companyGetter      companyBehavior.CompanyGetter
	userGetter         companyBehavior.UserGetter
	roleGetter         roleBehavior.RoleGetter
	roleAssigner       companyBehavior.RoleAssigner
}

func NewAssignControl(
	database db.TransactionalDatabase,
	authorization *authorization.Authorization,
	companyGetter companyBehavior.CompanyGetter,
	userGetter companyBehavior.UserGetter,
	roleGetter roleBehavior.RoleGetter,
	roleAssigner companyBehavior.RoleAssigner,
) *AssignControl {
	return &AssignControl{
		TransactionalDatabase: database,
		authorization:	       authorization,
		companyGetter:         companyGetter,
		userGetter:            userGetter,
		roleGetter:            roleGetter,
		roleAssigner:          roleAssigner,
	}
}

func CreateAssignControl() (*AssignControl, error) {
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

	companyGetter := companyBehavior.NewCompanyGet(database)
	userGetter := companyBehavior.NewUserGet(local, database)
	roleGetter := roleBehavior.NewRoleGet(database)
	roleAssigner := companyBehavior.NewRoleAssign(local, database)

	return NewAssignControl(database, authorization, companyGetter, userGetter, roleGetter, roleAssigner), nil
}

var assignCompanyPermission = shelterRole.NewRequirePermission(true, true, true, false)

var AssignExecute = utility.Transact(func(control *AssignControl, entry entryCompanyUser.CompanyUserAssignRequest, authentic *pkgUser.Authentic) (*entryCompanyUser.CompanyUserResponse, error) {

	if err := control.authorization.Authorize(assignCompanyPermission, authentic); err != nil {
		return nil, err
	}

	company, err := control.companyGetter.Execute(entry)
	if err != nil {
		return nil, err
	}
	if company == nil {
		keys := map[string]string{"identifier": entry.CompanyGet.Identifier}
		return nil, essence.NewNotFoundError("company", keys, "company not found")
	}

	userAuthenticWouldAssign, err := control.userGetter.Execute(entry, *company)
	if err != nil {
		return nil, err
	}
	if userAuthenticWouldAssign == nil {
		keys := map[string]string{"identifier": entry.RoleAssign.UserIdentifier}
		return nil, essence.NewNotFoundError("user", keys, "user not found")
	}

	role, err := control.roleGetter.Execute(entry)
	if err != nil {
		return nil, err
	}
	if role == nil {
		keys := map[string]string{"label": entry.RoleAssign.RoleLabel}
		return nil, essence.NewNotFoundError("role", keys, "role not found")
	}

	userAuthentic, err := control.roleAssigner.Execute(*company, userAuthenticWouldAssign, *role)
	if err != nil {
		return nil, err
	}

	return entryCompanyUser.FromShelterUserAuthenticToGetResponse(userAuthentic), nil
})
