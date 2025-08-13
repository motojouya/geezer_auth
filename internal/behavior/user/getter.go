package user

import (
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type UserGetterDB interface {
	userQuery.GetUserAuthenticQuery
}

type UserGetter interface {
	Execute(userIdentifier pkgText.Identifier) (*shelterUser.UserAuthentic, error)
}

type UserGet struct {
	local localPkg.Localer
	db    UserGetterDB
}

func NewUserGet(local localPkg.Localer, db UserGetterDB) *UserGet {
	return &UserGet{
		local: local,
		db:    db,
	}
}

func (getter UserGet) Execute(userIdentifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
	now := getter.local.GetNow()

	dbUserAuthentic, err := getter.db.GetUserAuthentic(string(userIdentifier), now)
	if err != nil {
		return nil, err
	}

	if dbUserAuthentic == nil {
		keys := map[string]string{"identifier": string(userIdentifier)}
		return nil, essence.NewNotFoundError("user", keys, "user not found")
	}

	return dbUserAuthentic.ToShelterUserAuthentic()
}
