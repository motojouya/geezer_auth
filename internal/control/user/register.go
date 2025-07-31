package user

import (
	"github.com/motojouya/geezer_auth/internal/db"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	"github.com/motojouya/geezer_auth/internal/io"
	configSilo "github.com/motojouya/geezer_auth/internal/silo/config"
	userSilo "github.com/motojouya/geezer_auth/internal/silo/user"
	pkgUser "github.com/motojouya/geezer_auth/pkg/core/user"
	"github.com/motojouya/geezer_auth/internal/control/utility"
)

type RegisterControl struct {
	db.TransactionalDatabase
	userCreator        *userSilo.UserCreator
	emailSetter        *userSilo.EmailSetter
	passwordSetter     *userSilo.PasswordSetter
	refreshTokenIssuer *userSilo.RefreshTokenIssuer
	accessTokenIssuer  *userSilo.AccessTokenIssuer
}

func NewRegisterControl(
	database db.TransactionalDatabase,
	userCreator *userSilo.UserCreator,
	emailSetter *userSilo.EmailSetter,
	passwordSetter *userSilo.PasswordSetter,
	refreshTokenIssuer *userSilo.RefreshTokenIssuer,
	accessTokenIssuer *userSilo.AccessTokenIssuer,
) *RegisterControl {
	return &RegisterControl{
		TransactionalDatabase: database,
		userCreator:           userCreator,
		emailSetter:           emailSetter,
		passwordSetter:        passwordSetter,
		refreshTokenIssuer:    refreshTokenIssuer,
		accessTokenIssuer:     accessTokenIssuer,
	}
}

func CreateRegisterControl() (*RegisterControl, error) {
	var local = io.CreateLocal()
	var env = io.CreateEnvironment()

	database, err := configSilo.NewDatabaseLoader(env).LoadDatabase()
	if err != nil {
		return nil, err
	}

	jwtHandler, err := configSilo.NewJwtHandlerLoader(env).LoadJwtHandler()
	if err != nil {
		return nil, err
	}

	userCreator := userSilo.NewUserCreator(local, database)
	emailSetter := userSilo.NewEmailSetter(local, database)
	passwordSetter := userSilo.NewPasswordSetter(local, database)
	refreshTokenIssuer := userSilo.NewRefreshTokenIssuer(local, database)
	accessTokenIssuer := userSilo.NewAccessTokenIssuer(local, database, jwtHandler)

	return NewRegisterControl(database, userCreator, emailSetter, passwordSetter, refreshTokenIssuer, accessTokenIssuer), nil
}

var RegisterExecute = utility.Transact(func(control *RegisterControl, entry entryUser.UserRegisterRequest, _ *pkgUser.User) (*entryUser.UserRegisterResponse, error) {

	userAuthentic, err := control.userCreator.Execute(entry)
	if err != nil {
		return nil, err
	}

	if err = control.passwordSetter.Execute(entry, userAuthentic); err != nil {
		return nil, err
	}

	if err = control.emailSetter.Execute(entry, userAuthentic); err != nil {
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

	return entryUser.FromCoreUserAuthenticToRegisterResponse(userAuthentic, refreshToken, accessToken), nil
})
