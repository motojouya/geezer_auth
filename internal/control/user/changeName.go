package user

import (
	authBehavior "github.com/motojouya/geezer_auth/internal/behavior/authorization"
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	userBehavior "github.com/motojouya/geezer_auth/internal/behavior/user"
	"github.com/motojouya/geezer_auth/internal/control/utility"
	"github.com/motojouya/geezer_auth/internal/db"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/authorization"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
)

type ChangeNameControl struct {
	db.TransactionalDatabase
	authorization     *authorization.Authorization
	userGetter        userBehavior.UserGetter
	nameChanger       userBehavior.NameChanger
	accessTokenIssuer userBehavior.AccessTokenIssuer
}

func NewChangeNameControl(
	database db.TransactionalDatabase,
	authorization *authorization.Authorization,
	userGetter userBehavior.UserGetter,
	nameChanger userBehavior.NameChanger,
	accessTokenIssuer userBehavior.AccessTokenIssuer,
) *ChangeNameControl {
	return &ChangeNameControl{
		TransactionalDatabase: database,
		authorization:         authorization,
		userGetter:            userGetter,
		nameChanger:           nameChanger,
		accessTokenIssuer:     accessTokenIssuer,
	}
}

func CreateChangeNameControl() (*ChangeNameControl, error) {
	var local = localPkg.CreateLocal()
	var env = localPkg.CreateEnvironment()

	database, err := configBehavior.NewDatabaseGet(env).GetDatabase()
	if err != nil {
		return nil, err
	}

	jwtHandler, err := configBehavior.NewJwtHandlerGet(env).GetJwtHandler()
	if err != nil {
		return nil, err
	}

	authorization, err := authBehavior.NewAuthorizationGet(database).GetAuthorization()
	if err != nil {
		return nil, err
	}

	userGetter := userBehavior.NewUserGet(local, database)
	nameChanger := userBehavior.NewNameChange(local, database)
	accessTokenIssuer := userBehavior.NewAccessTokenIssue(local, database, jwtHandler)

	return NewChangeNameControl(database, authorization, userGetter, nameChanger, accessTokenIssuer), nil
}

var changeNamePermission = shelterRole.NewRequirePermission(true, false, false, false)

var ChangeNameExecute = utility.Transact(func(control *ChangeNameControl, entry entryUser.UserChangeNameRequest, authentic *pkgUser.Authentic) (*entryUser.UserUpdateResponse, error) {

	if err := control.authorization.Authorize(changeNamePermission, authentic); err != nil {
		return nil, err
	}

	userAuthentic, err := control.userGetter.Execute(authentic.User.Identifier)
	if err != nil {
		return nil, err
	}

	userAuthentic, err = control.nameChanger.Execute(entry, userAuthentic)
	if err != nil {
		return nil, err
	}

	accessToken, err := control.accessTokenIssuer.Execute(userAuthentic)
	if err != nil {
		return nil, err
	}

	return entryUser.FromShelterUserAuthenticToUpdateResponse(userAuthentic, accessToken), nil
})
