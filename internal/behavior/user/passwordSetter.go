package user

import (
	"github.com/go-gorp/gorp"
	commandQuery "github.com/motojouya/geezer_auth/internal/db/query/command"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
)

type PasswordSetterDB interface {
	gorp.SqlExecutor
	commandQuery.AddPasswordQuery
}

type PasswordSetter interface {
	Execute(entry entryUser.PasswordGetter, userAuthentic *shelterUser.UserAuthentic) error
}

type PasswordSet struct {
	local localPkg.Localer
	db    PasswordSetterDB
}

func NewPasswordSet(local localPkg.Localer, db PasswordSetterDB) *PasswordSet {
	return &PasswordSet{
		local: local,
		db:    db,
	}
}

func (setter PasswordSet) Execute(entry entryUser.PasswordGetter, userAuthentic *shelterUser.UserAuthentic) error {
	now := setter.local.GetNow()

	password, err := entry.GetPassword()
	if err != nil {
		return err
	}

	hashedPassword, err := shelterText.HashPassword(password)
	if err != nil {
		return err
	}

	userPassword := shelterUser.CreateUserPassword(userAuthentic.GetUser(), hashedPassword, now)

	dbUserPassword := dbUser.FromCoreUserPassword(userPassword)

	if err = setter.db.Insert(dbUserPassword); err != nil {
		return err
	}

	return nil
}
