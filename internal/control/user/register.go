package user

import (
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	userBehavior "github.com/motojouya/geezer_auth/internal/behavior/user"
	"github.com/motojouya/geezer_auth/internal/control/utility"
	"github.com/motojouya/geezer_auth/internal/db"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
)

type RegisterControl struct {
	db.TransactionalDatabase
	userCreator        userBehavior.UserCreator
	emailSetter        userBehavior.EmailSetter
	passwordSetter     userBehavior.PasswordSetter
	refreshTokenIssuer userBehavior.RefreshTokenIssuer
	accessTokenIssuer  userBehavior.AccessTokenIssuer
}

func NewRegisterControl(
	database db.TransactionalDatabase,
	userCreator userBehavior.UserCreator,
	emailSetter userBehavior.EmailSetter,
	passwordSetter userBehavior.PasswordSetter,
	refreshTokenIssuer userBehavior.RefreshTokenIssuer,
	accessTokenIssuer userBehavior.AccessTokenIssuer,
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

	userCreator := userBehavior.NewUserCreate(local, database)
	emailSetter := userBehavior.NewEmailSet(local, database)
	passwordSetter := userBehavior.NewPasswordSet(local, database)
	refreshTokenIssuer := userBehavior.NewRefreshTokenIssue(local, database)
	accessTokenIssuer := userBehavior.NewAccessTokenIssue(local, database, jwtHandler)

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

	return entryUser.FromShelterUserAuthenticToRegisterResponse(userAuthentic, refreshToken, accessToken), nil
})
