package role

import (
	authBehavior "github.com/motojouya/geezer_auth/internal/behavior/authorization"
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	roleBehavior "github.com/motojouya/geezer_auth/internal/behavior/role"
	entryCommon "github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	entryCompanyUser "github.com/motojouya/geezer_auth/internal/entry/transfer/companyUser"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/authorization"
	essence "github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
)

type GetRoleControl struct {
	essence.Closable
	authorization *authorization.Authorization
	roleGetter    roleBehavior.AllRoleGetter
}

func NewGetRoleControl(
	database essence.Closable,
	authorization *authorization.Authorization,
	roleGetter    roleBehavior.AllRoleGetter,
) *GetRoleControl {
	return &GetRoleControl{
		Closable:      database,
		authorization: authorization,
		roleGetter:    roleGetter,
	}
}

func CreateGetRoleControl() (*GetRoleControl, error) {
	var env = localPkg.CreateEnvironment()

	database, err := configBehavior.NewDatabaseGet(env).GetDatabase()
	if err != nil {
		return nil, err
	}

	authorization, err := authBehavior.NewAuthorizationGet(database).GetAuthorization()
	if err != nil {
		return nil, err
	}

	roleGetter := roleBehavior.NewRoleGet(database)

	return NewGetRoleControl(database, authorization, roleGetter), nil
}

var getRolePermission = shelterRole.NewRequirePermission(true, false, false, false)

func GetRoleExecute(control *GetRoleControl, entry entryCommon.Empty, authentic *pkgUser.Authentic) (*entryCompanyUser.RoleGetResponse, error) {

	if err := control.authorization.Authorize(getRolePermission, authentic); err != nil {
		return nil, err
	}

	roles, err := control.roleGetter.Execute()
	if err != nil {
		return nil, err
	}

	return entryCompanyUser.FromShelterRoles(roles), nil
}
