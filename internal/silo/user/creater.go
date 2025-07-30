package user

import (
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/core/essence"
	coreText "github.com/motojouya/geezer_auth/internal/core/text"
	coreUser "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/db"
	commandQuery "github.com/motojouya/geezer_auth/internal/db/query/command"
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	"github.com/motojouya/geezer_auth/internal/io"
	"github.com/motojouya/geezer_auth/internal/service"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type UserCreatorDB interface {
	gorp.SqlExecutor
	userQuery.GetUserQuery
	userQuery.GetUserAuthenticQuery
}

type UserCreator struct {
	local io.Local
	db    UserCreatorDB
}

func NewUserCreator(local io.Local, db UserCreatorDB) *UserCreator {
	return &RegisterControl{
		local: local,
		db:    database,
	}
}

func createUserIdentifier(local io.Local) func() (pkgText.Identifier, error) {
	return func() (pkgText.Identifier, error) {
		var ramdomString = local.GenerateRamdomString(pkgText.IdentifierLength, pkgText.IdentifierChar)
		var identifier, err = coreUser.CreateUserIdentifier(ramdomString)
		if err != nil {
			return pkgText.Identifier(""), err
		}
		return identifier, nil
	}
}

func checkUserIdentifier(userCreatorDB UserCreatorDB) func(pkgText.Identifier) (bool, error) {
	return func(identifier pkgText.Identifier) (bool, error) {
		var user, err = userCreatorDB.GetUser(string(identifier))
		if err != nil {
			return false, err
		}
		return user == nil, nil
	}
}

type UserGetter interface {
	ToCoreUser(pkgText.Identifier, time.Time) (coreUser.UnsavedUser, error)
}

func (creator UserCreator) Create(entry UserGetter) (*coreUser.UserAuthentic, error) {
	now := creator.local.GetNow()

	userIdentifier, err := coreText.GetString(createUserIdentifier(creator.local), checkUserIdentifier(creator.local), 10)
	if err != nil {
		return coreUser.User{}, err
	}

	unsavedUser, err := entry.ToCoreUser(userIdentifier, now)
	if err != nil {
		return coreUser.User{}, err
	}

	var dbUserValue = dbUser.FromCoreUser(unsavedUser)

	if err = creator.local.Insert(&dbUserValue); err != nil {
		return coreUser.User{}, err
	}

	dbUserAuthentic, err := creator.local.GetUserAuthentic(dbUserValue.Identifier, now)
	if err != nil {
		return nil, err
	}

	if dbUserAuthentic == nil {
		keys := map[string]string{"identifier": string(savedUser.Identifier)}
		err = essence.NewNotFoundError("user", keys, "user not found")
		return nil, err
	}

	return dbUserAuthentic.ToCoreUserAuthentic()

}
