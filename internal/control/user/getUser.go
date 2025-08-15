package user

import (
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	authBehavior "github.com/motojouya/geezer_auth/internal/behavior/authorization"
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	userBehavior "github.com/motojouya/geezer_auth/internal/behavior/user"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	entryCommon "github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
	"github.com/motojouya/geezer_auth/internal/shelter/authorization"
	essence "github.com/motojouya/geezer_auth/internal/shelter/essence"
)

type GetUserControl struct {
	essence.Closable
	authorization *authorization.Authorization
	userGetter    userBehavior.UserGetter
}

func NewGetUserControl(
	database essence.Closable,
	authorization *authorization.Authorization,
	userGetter userBehavior.UserGetter,
) *GetUserControl {
	return &GetUserControl{
		Closable:      database,
		authorization: authorization,
		userGetter:    userGetter,
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

	userGetter := userBehavior.NewUserGet(local, database)

	return NewGetUserControl(database, authorization, userGetter), nil
}

var getUserPermission = shelterRole.NewRequirePermission(true, false, false, false)

func GetUserExecute(control *GetUserControl, entry entryCommon.Empty, authentic *pkgUser.Authentic) (*entryUser.UserGetResponse, error) {

	if err := control.authorization.Authorize(getUserPermission, authentic); err != nil {
		return nil, err
	}

	userAuthentic, err := control.userGetter.Execute(authentic.User.Identifier)
	if err != nil {
		return nil, err
	}

	return entryUser.FromShelterUserAuthenticToGetResponse(userAuthentic), nil
}
