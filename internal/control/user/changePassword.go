package user

import (
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	authBehavior "github.com/motojouya/geezer_auth/internal/behavior/authorization"
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	userBehavior "github.com/motojouya/geezer_auth/internal/behavior/user"
	"github.com/motojouya/geezer_auth/internal/control/utility"
	"github.com/motojouya/geezer_auth/internal/db"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
	"github.com/motojouya/geezer_auth/internal/shelter/authorization"
)

type ChangePasswordControl struct {
	db.TransactionalDatabase
	authorization     *authorization.Authorization
	userGetter        userBehavior.UserGetter
	passwordSetter    userBehavior.PasswordSetter
}

func NewChangePasswordControl(
	database db.TransactionalDatabase,
	authorization *authorization.Authorization,
	userGetter userBehavior.UserGetter,
	passwordSetter userBehavior.PasswordSetter,
) *ChangePasswordControl {
	return &ChangePasswordControl{
		TransactionalDatabase: database,
		authorization:         authorization,
		userGetter:            userGetter,
		passwordSetter:        passwordSetter,
	}
}

func CreateChangePasswordControl() (*ChangePasswordControl, error) {
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
	passwordSetter := userBehavior.NewPasswordSet(local, database)

	return NewChangePasswordControl(database, authorization, userGetter, passwordSetter), nil
}

var changePasswordPermission = shelterRole.NewRequirePermission(true, false, false, false)

var ChangePasswordExecute = utility.Transact(func(control *ChangePasswordControl, entry entryUser.UserChangePasswordRequest, authentic *pkgUser.Authentic) (*entryUser.UserGetResponse, error) {

	if err := control.authorization.Authorize(changePasswordPermission, authentic); err != nil {
		return nil, err
	}

	userAuthentic, err := control.userGetter.Execute(authentic.User.Identifier)
	if err != nil {
		return nil, err
	}

	if err = control.passwordSetter.Execute(entry, userAuthentic); err != nil {
		return nil, err
	}

	return entryUser.FromShelterUserAuthenticToGetResponse(userAuthentic), nil
})
