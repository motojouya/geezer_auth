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

type ChangeEmailControl struct {
	db.TransactionalDatabase
	authorization *authorization.Authorization
	userGetter    userBehavior.UserGetter
	emailSetter   userBehavior.EmailSetter
}

func NewChangeEmailControl(
	database db.TransactionalDatabase,
	authorization *authorization.Authorization,
	userGetter userBehavior.UserGetter,
	emailSetter userBehavior.EmailSetter,
) *ChangeEmailControl {
	return &ChangeEmailControl{
		TransactionalDatabase: database,
		authorization:         authorization,
		userGetter:            userGetter,
		emailSetter:           emailSetter,
	}
}

func CreateChangeEmailControl() (*ChangeEmailControl, error) {
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
	emailSetter := userBehavior.NewEmailSet(local, database)

	return NewChangeEmailControl(database, authorization, userGetter, emailSetter), nil
}

var changeEmailPermission = shelterRole.NewRequirePermission(true, false, false, false)

var ChangeEmailExecute = utility.Transact(func(control *ChangeEmailControl, entry entryUser.UserChangeEmailRequest, authentic *pkgUser.Authentic) (*entryUser.UserGetResponse, error) {

	if err := control.authorization.Authorize(changeEmailPermission, authentic); err != nil {
		return nil, err
	}

	userAuthentic, err := control.userGetter.Execute(authentic.User.Identifier)
	if err != nil {
		return nil, err
	}

	if err = control.emailSetter.Execute(entry, userAuthentic); err != nil {
		return nil, err
	}

	return entryUser.FromShelterUserAuthenticToGetResponse(userAuthentic), nil
})
