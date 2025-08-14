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

type VerifyEmailControl struct {
	db.TransactionalDatabase
	Authorization     *authorization.Authorization
	userGetter        userBehavior.UserGetter
	emailVerifier     userBehavior.EmailVerifier
	accessTokenIssuer userBehavior.AccessTokenIssuer
}

func NewVerifyEmailControl(
	database db.TransactionalDatabase,
	authorization *authorization.Authorization,
	userGetter userBehavior.UserGetter,
	emailVerifier userBehavior.EmailVerifier,
	accessTokenIssuer userBehavior.AccessTokenIssuer,
) *VerifyEmailControl {
	return &VerifyEmailControl{
		TransactionalDatabase: database,
		Authorization:         authorization,
		userGetter:            userGetter,
		emailVerifier:         emailVerifier,
		accessTokenIssuer:     accessTokenIssuer,
	}
}

func CreateVerifyEmailControl() (*VerifyEmailControl, error) {
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
	emailVerifier := userBehavior.NewEmailVerify(local, database)
	accessTokenIssuer := userBehavior.NewAccessTokenIssue(local, database, jwtHandler)

	return NewVerifyEmailControl(database, authorization, userGetter, emailVerifier, accessTokenIssuer), nil
}

var permission = shelterRole.NewRequirePermission(true, false, false, false)

var EmailVerifyExecute = utility.Transact(func(control *VerifyEmailControl, entry entryUser.UserVerifyEmailRequest, authentic *pkgUser.Authentic) (*entryUser.UserUpdateResponse, error) {

	if err := control.Authorization.Authorize(permission, authentic); err != nil {
		return nil, err
	}

	userAuthentic, err := control.userGetter.Execute(authentic.User.Identifier)
	if err != nil {
		return nil, err
	}

	userAuthentic, err = control.emailVerifier.Execute(entry, userAuthentic)
	if err != nil {
		return nil, err
	}

	accessToken, err := control.accessTokenIssuer.Execute(userAuthentic)
	if err != nil {
		return nil, err
	}

	return entryUser.FromShelterUserAuthenticToUpdateResponse(userAuthentic, accessToken), nil
})
