package companyUser

import (
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	authBehavior "github.com/motojouya/geezer_auth/internal/behavior/authorization"
	userBehavior "github.com/motojouya/geezer_auth/internal/behavior/user"
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

type AcceptControl struct {
	db.TransactionalDatabase
	authorization      *authorization.Authorization
	userGetter         userBehavior.UserGetter
	companyGetter      companyBehavior.CompanyGetter
	inviteTokenChecker companyBehavior.InviteTokenChecker
	roleAssigner       companyBehavior.RoleAssigner
}

func NewAcceptControl(
	database db.TransactionalDatabase,
	authorization *authorization.Authorization,
	userGetter userBehavior.UserGetter,
	companyGetter companyBehavior.CompanyGetter,
	inviteTokenChecker companyBehavior.InviteTokenChecker,
	roleAssigner companyBehavior.RoleAssigner,
) *AcceptControl {
	return &AcceptControl{
		TransactionalDatabase: database,
		authorization:	       authorization,
		userGetter:            userGetter,
		companyGetter:         companyGetter,
		inviteTokenChecker:    inviteTokenChecker,
		roleAssigner:          roleAssigner,
	}
}

func CreateAcceptControl() (*AcceptControl, error) {
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

	userGetter := userBehavior.NewUserGet(local, database)
	companyGetter := companyBehavior.NewCompanyGet(database)
	inviteTokenCheck := companyBehavior.NewInviteTokenCheck(local, database)
	roleAssigner := companyBehavior.NewRoleAssign(local, database)

	return NewAcceptControl(database, authorization, userGetter, companyGetter, inviteTokenCheck, roleAssigner), nil
}

var acceptCompanyPermission = shelterRole.NewRequirePermission(true, false, false, false)

var AcceptExecute = utility.Transact(func(control *AcceptControl, entry entryCompanyUser.CompanyUserAcceptRequest, authentic *pkgUser.Authentic) (*entryCompanyUser.CompanyUserResponse, error) {

	if err := control.authorization.Authorize(acceptCompanyPermission, authentic); err != nil {
		return nil, err
	}

	userAuthenticWouldAssign, err := control.userGetter.Execute(authentic.User.Identifier)
	if err != nil {
		return nil, err
	}
	if userAuthenticWouldAssign == nil {
		keys := map[string]string{"identifier": string(authentic.User.Identifier)}
		return nil, essence.NewNotFoundError("user", keys, "user not found")
	}

	company, err := control.companyGetter.Execute(entry)
	if err != nil {
		return nil, err
	}
	if company == nil {
		keys := map[string]string{"identifier": entry.CompanyGet.Identifier}
		return nil, essence.NewNotFoundError("company", keys, "company not found")
	}

	role, err := control.inviteTokenChecker.Execute(entry, *company)
	if err != nil {
		return nil, err
	}

	userAuthentic, err := control.roleAssigner.Execute(*company, userAuthenticWouldAssign, role)
	if err != nil {
		return nil, err
	}

	return entryCompanyUser.FromShelterUserAuthenticToGetResponse(userAuthentic), nil
})
