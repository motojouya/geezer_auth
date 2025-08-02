package user

import (
	"github.com/motojouya/geezer_auth/internal/control/utility"
	"github.com/motojouya/geezer_auth/internal/db"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	"github.com/motojouya/geezer_auth/internal/io"
	configSilo "github.com/motojouya/geezer_auth/internal/silo/config"
	userSilo "github.com/motojouya/geezer_auth/internal/silo/user"
	pkgUser "github.com/motojouya/geezer_auth/pkg/core/user"
)

type RegisterControl struct {
	db.TransactionalDatabase
	userCreator        userSilo.UserCreator
	emailSetter        userSilo.EmailSetter
	passwordSetter     userSilo.PasswordSetter
	refreshTokenIssuer userSilo.RefreshTokenIssuer
	accessTokenIssuer  userSilo.AccessTokenIssuer
}

func NewRegisterControl(
	database db.TransactionalDatabase,
	userCreator userSilo.UserCreator,
	emailSetter userSilo.EmailSetter,
	passwordSetter userSilo.PasswordSetter,
	refreshTokenIssuer userSilo.RefreshTokenIssuer,
	accessTokenIssuer userSilo.AccessTokenIssuer,
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

	database, err := configSilo.NewDatabaseGet(env).GetDatabase()
	if err != nil {
		return nil, err
	}

	jwtHandler, err := configSilo.NewJwtHandlerGet(env).GetJwtHandler()
	if err != nil {
		return nil, err
	}

	userCreator := userSilo.NewUserCreate(local, database)
	emailSetter := userSilo.NewEmailSet(local, database)
	passwordSetter := userSilo.NewPasswordSet(local, database)
	refreshTokenIssuer := userSilo.NewRefreshTokenIssue(local, database)
	accessTokenIssuer := userSilo.NewAccessTokenIssue(local, database, jwtHandler)

	return NewRegisterControl(database, userCreator, emailSetter, passwordSetter, refreshTokenIssuer, accessTokenIssuer), nil
}

var RegisterExecute = utility.Transact(func(control *RegisterControl, entry entryUser.UserRegisterRequest, _ *pkgUser.Authentic) (*entryUser.UserRegisterResponse, error) {

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
