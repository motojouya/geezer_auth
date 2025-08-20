package user

import (
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	entryAuth "github.com/motojouya/geezer_auth/internal/entry/transfer/auth"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
)

type PasswordCheckerDB interface {
	userQuery.GetUserPasswordQuery
}

type PasswordChecker interface {
	Execute(entry entryAuth.AuthLoginner) error
}

type PasswordCheck struct {
	local localPkg.Localer
	db    PasswordCheckerDB
}

func NewPasswordCheck(local localPkg.Localer, db PasswordSetterDB) *PasswordCheck {
	return &PasswordCheck{
		local: local,
		db:    db,
	}
}

func (checker PasswordCheck) Execute(entry entryAuth.AuthLoginner) error {
	now := checker.local.GetNow()

	identifier, err := entry.GetIdentifier()
	if err != nil {
		return err
	}

	if identifier != nil {
	}

	password, err := entry.GetPassword()
	if err != nil {
		return err
	}

	hashedPassword, err := shelterText.HashPassword(password)
	if err != nil {
		return err
	}

	userPassword := shelterUser.CreateUserPassword(userAuthentic.GetUser(), hashedPassword, now)

	dbUserPassword := dbUser.FromShelterUserPassword(userPassword)

	_, err = checker.db.GetUserPassword(dbUserPassword, now)
	if err != nil {
		return err
	}

	return nil
}
