package user

import (
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type UserCreatorDB interface {
	gorp.SqlExecutor
	userQuery.GetUserQuery
	userQuery.GetUserAuthenticQuery
}

type UserCreator interface {
	Execute(entry entryUser.UserGetter) (*shelterUser.UserAuthentic, error)
}

type UserCreate struct {
	local localPkg.Localer
	db    UserCreatorDB
}

func NewUserCreate(local localPkg.Localer, db UserCreatorDB) *UserCreate {
	return &UserCreate{
		local: local,
		db:    db,
	}
}

func createUserIdentifier(local localPkg.Localer) func() (pkgText.Identifier, error) {
	return func() (pkgText.Identifier, error) {
		var ramdomString = local.GenerateRamdomString(pkgText.IdentifierLength, pkgText.IdentifierChar)
		var identifier, err = shelterUser.CreateUserIdentifier(ramdomString)
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

func (creator UserCreate) Execute(entry entryUser.UserGetter) (*shelterUser.UserAuthentic, error) {
	now := creator.local.GetNow()

	userIdentifier, err := shelterText.GetString(createUserIdentifier(creator.local), checkUserIdentifier(creator.db), 10)
	if err != nil {
		return nil, err
	}

	unsavedUser, err := entry.ToCoreUser(userIdentifier, now)
	if err != nil {
		return nil, err
	}

	var dbUserValue = dbUser.FromCoreUser(unsavedUser)

	if err = creator.db.Insert(&dbUserValue); err != nil {
		return nil, err
	}

	dbUserAuthentic, err := creator.db.GetUserAuthentic(dbUserValue.Identifier, now)
	if err != nil {
		return nil, err
	}

	if dbUserAuthentic == nil {
		keys := map[string]string{"identifier": dbUserValue.Identifier}
		err = essence.NewNotFoundError("user", keys, "user not found")
		return nil, err
	}

	return dbUserAuthentic.ToCoreUserAuthentic()
}
