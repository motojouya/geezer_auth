package user

import (
	"github.com/go-gorp/gorp"
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
)

type NameChangerDB interface {
	gorp.SqlExecutor
	userQuery.GetUserAuthenticQuery
}

type NameChanger interface {
	Execute(entry entryUser.UserApplyer, userAuthentic *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error)
}

type NameChange struct {
	local localPkg.Localer
	db    NameChangerDB
}

func NewNameChange(local localPkg.Localer, db NameChangerDB) *NameChange {
	return &NameChange{
		local: local,
		db:    db,
	}
}

func (changer NameChange) Execute(entry entryUser.UserApplyer, userAuthentic *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error) {
	now := changer.local.GetNow()

	changeApplyedUser, err := entry.ApplyShelterUser(userAuthentic.GetUser(), now)
	if err != nil {
		return nil, err
	}

	var dbUserValue = dbUser.FromShelterUser(changeApplyedUser)

	if _, err = changer.db.Update(&dbUserValue); err != nil {
		return nil, err
	}

	dbUserAuthentic, err := changer.db.GetUserAuthentic(dbUserValue.Identifier, now)
	if err != nil {
		return nil, err
	}

	if dbUserAuthentic == nil {
		keys := map[string]string{"identifier": dbUserValue.Identifier}
		return nil, essence.NewNotFoundError("user", keys, "user not found")
	}

	return dbUserAuthentic.ToShelterUserAuthentic()
}
