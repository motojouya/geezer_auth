package user

import (
	"github.com/go-gorp/gorp"
	coreText "github.com/motojouya/geezer_auth/internal/core/text"
	coreUser "github.com/motojouya/geezer_auth/internal/core/user"
	commandQuery "github.com/motojouya/geezer_auth/internal/db/query/command"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	"github.com/motojouya/geezer_auth/internal/io"
)

type PasswordSetterDB interface {
	gorp.SqlExecutor
	commandQuery.AddPasswordQuery
}

type PasswordSetter interface {
	Execute(entry entryUser.PasswordGetter, userAuthentic *coreUser.UserAuthentic) error
}

type PasswordSet struct {
	local io.Local
	db    PasswordSetterDB
}

func NewPasswordSet(local io.Local, db PasswordSetterDB) *PasswordSet {
	return &PasswordSet{
		local: local,
		db:    db,
	}
}

func (setter PasswordSet) Execute(entry entryUser.PasswordGetter, userAuthentic *coreUser.UserAuthentic) error {
	now := setter.local.GetNow()

	password, err := entry.GetPassword()
	if err != nil {
		return err
	}

	hashedPassword, err := coreText.HashPassword(password)
	if err != nil {
		return err
	}

	userPassword := coreUser.CreateUserPassword(userAuthentic.GetUser(), hashedPassword, now)

	dbUserPassword := dbUser.FromCoreUserPassword(userPassword)

	if err = setter.db.Insert(dbUserPassword); err != nil {
		return err
	}

	return nil
}
