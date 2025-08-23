package user

import (
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	userBehavior "github.com/motojouya/geezer_auth/internal/behavior/user"
	"github.com/motojouya/geezer_auth/internal/control/utility"
	"github.com/motojouya/geezer_auth/internal/db"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	entryAuth "github.com/motojouya/geezer_auth/internal/entry/transfer/auth"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
)

type LoginControl struct {
	db.TransactionalDatabase
	userGetter         userBehavior.UserGetter
	passwordChecker    userBehavior.PasswordChecker
	refreshTokenIssuer userBehavior.RefreshTokenIssuer
	accessTokenIssuer  userBehavior.AccessTokenIssuer
}

func NewLoginControl(
	database db.TransactionalDatabase,
	userGetter userBehavior.UserGetter,
	passwordChecker userBehavior.PasswordChecker,
	refreshTokenIssuer userBehavior.RefreshTokenIssuer,
	accessTokenIssuer userBehavior.AccessTokenIssuer,
) *LoginControl {
	return &LoginControl{
		TransactionalDatabase: database,
		userGetter:            userGetter,
		passwordChecker:       passwordChecker,
		refreshTokenIssuer:    refreshTokenIssuer,
		accessTokenIssuer:     accessTokenIssuer,
	}
}

func CreateLoginControl() (*LoginControl, error) {
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

	userGetter := userBehavior.NewUserGet(local, database)
	passwordChecker := userBehavior.NewPasswordCheck(database)
	refreshTokenIssuer := userBehavior.NewRefreshTokenIssue(local, database)
	accessTokenIssuer := userBehavior.NewAccessTokenIssue(local, database, jwtHandler)

	return NewLoginControl(database, userGetter, passwordChecker, refreshTokenIssuer, accessTokenIssuer), nil
}

var LoginExecute = utility.Transact(func(control *LoginControl, entry entryAuth.AuthLoginRequest, _ *pkgUser.Authentic) (*entryAuth.AuthLoginResponse, error) {

	identifier, err := control.passwordChecker.Execute(entry)
	if err != nil {
		return nil, err
	}

	userAuthentic, err := control.userGetter.Execute(identifier)
	if err != nil {
		return nil, err
	}

	refreshToken, err := control.refreshTokenIssuer.Execute(userAuthentic)
	if err != nil {
		return nil, err
	}

	accessToken, err := control.accessTokenIssuer.Execute(userAuthentic)
	if err != nil {
		return nil, err
	}

	return entryAuth.FromShelterUserAuthenticToLoginResponse(userAuthentic, refreshToken, accessToken), nil
})
