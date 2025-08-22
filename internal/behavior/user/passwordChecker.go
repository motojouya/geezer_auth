package user

import (
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	entryAuth "github.com/motojouya/geezer_auth/internal/entry/transfer/auth"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
)

type PasswordCheckerDB interface {
	userQuery.GetUserPasswordQuery
	userQuery.GetUserPasswordOfEmailQuery
}

type PasswordChecker interface {
	Execute(entry entryAuth.AuthLoginner) error
}

type PasswordCheck struct {
	db PasswordCheckerDB
}

func NewPasswordCheck(db PasswordCheckerDB) *PasswordCheck {
	return &PasswordCheck{
		db: db,
	}
}

func (checker PasswordCheck) Execute(entry entryAuth.AuthLoginner) (pkgText.Identifier, error) {
	identifier, err := entry.GetIdentifier()
	if err != nil {
		return pkgText.Identifier(""), err
	}

	email, err := entry.GetEmailIdentifier()
	if err != nil {
		return pkgText.Identifier(""), err
	}

	password, err := entry.GetPassword()
	if err != nil {
		return pkgText.Identifier(""), err
	}

	var dbUserPassword *dbUser.UserPasswordFull
	if identifier != nil {
		dbUserPassword, err = checker.db.GetUserPassword(string(*identifier))
		if err != nil {
			return pkgText.Identifier(""), err
		}
	} else if email != nil {
		dbUserPassword, err = checker.db.GetUserPasswordOfEmail(string(*email))
		if err != nil {
			return pkgText.Identifier(""), err
		}
	} else {
		return pkgText.Identifier(""), essence.NewInvalidArgumentError("identifier", "", "identifier or email must be provided")
	}

	if dbUserPassword == nil {
		return pkgText.Identifier(""), essence.NewInvalidArgumentError("identifier", "", "user not found")
	}

	userPassword, err := dbUserPassword.ToShelterUserPassword()
	if err != nil {
		return pkgText.Identifier(""), err
	}

	err = shelterText.VerifyPassword(userPassword.Password, password)
	if err != nil {
		return pkgText.Identifier(""), err
	}

	return userPassword.User.Identifier, nil
}
