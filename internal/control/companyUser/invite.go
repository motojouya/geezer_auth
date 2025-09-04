package companyUser

import (
	authBehavior "github.com/motojouya/geezer_auth/internal/behavior/authorization"
	companyBehavior "github.com/motojouya/geezer_auth/internal/behavior/company"
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	roleBehavior "github.com/motojouya/geezer_auth/internal/behavior/role"
	"github.com/motojouya/geezer_auth/internal/control/utility"
	"github.com/motojouya/geezer_auth/internal/db"
	entryCompanyUser "github.com/motojouya/geezer_auth/internal/entry/transfer/companyUser"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/authorization"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
)

type InviteControl struct {
	db.TransactionalDatabase
	authorization     *authorization.Authorization
	companyGetter     companyBehavior.CompanyGetter
	roleGetter        roleBehavior.RoleGetter
	inviteTokenIssuer companyBehavior.InviteTokenIssuer
}

func NewInviteControl(
	database db.TransactionalDatabase,
	authorization *authorization.Authorization,
	companyGetter companyBehavior.CompanyGetter,
	roleGetter roleBehavior.RoleGetter,
	inviteTokenIssuer companyBehavior.InviteTokenIssuer,
) *InviteControl {
	return &InviteControl{
		TransactionalDatabase: database,
		authorization:         authorization,
		companyGetter:         companyGetter,
		roleGetter:            roleGetter,
		inviteTokenIssuer:     inviteTokenIssuer,
	}
}

func CreateInviteControl() (*InviteControl, error) {
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
	roleGetter := roleBehavior.NewRoleGet(database)
	inviteTokenIssuer := companyBehavior.NewInviteTokenIssue(local, database)

	return NewInviteControl(database, authorization, companyGetter, roleGetter, inviteTokenIssuer), nil
}

var inviteCompanyPermission = shelterRole.NewRequirePermission(true, true, true, false)

var InviteExecute = utility.Transact(func(control *InviteControl, entry entryCompanyUser.CompanyUserInviteRequest, authentic *pkgUser.Authentic) (entryCompanyUser.CompanyUserInviteResponse, error) {

	if err := control.authorization.Authorize(inviteCompanyPermission, authentic); err != nil {
		return entryCompanyUser.CompanyUserInviteResponse{}, err
	}

	company, err := control.companyGetter.Execute(entry)
	if err != nil {
		return entryCompanyUser.CompanyUserInviteResponse{}, err
	}
	if company == nil {
		keys := map[string]string{"identifier": entry.CompanyGet.Identifier}
		return entryCompanyUser.CompanyUserInviteResponse{}, essence.NewNotFoundError("company", keys, "company not found")
	}

	role, err := control.roleGetter.Execute(entry)
	if err != nil {
		return entryCompanyUser.CompanyUserInviteResponse{}, err
	}
	if role == nil {
		keys := map[string]string{"label": entry.RoleInvite.RoleLabel}
		return entryCompanyUser.CompanyUserInviteResponse{}, essence.NewNotFoundError("role", keys, "role not found")
	}

	token, err := control.inviteTokenIssuer.Execute(*company, *role)
	if err != nil {
		return entryCompanyUser.CompanyUserInviteResponse{}, err
	}

	return entryCompanyUser.FromToken(token), nil
})
