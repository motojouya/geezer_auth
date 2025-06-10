package user

import (
	text "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

const UserIdentifierPrefix = "US-"

type UnsavedUser struct {
	Identifier     text.Identifier
	ExposeEmailId  text.Email
	Name           text.Name
	BotFlag        bool
	RegisteredDate time.Time
	UpdateDate     time.Time
}

type User struct {
	PersistKey uint
	UnsavedUser
}

func CreateUserIdentifier(random string) (text.Identifier, error) {
	return text.CreateIdentifier(UserIdentifierPrefix, random)
}

func CreateUser(
	identifier text.Identifier,
	emailId text.Email,
	name text.Name,
	botFlag bool,
	registeredDate time.Time,
) UnsavedUser {
	return UnsavedUser{
		Identifier:     identifier,
		ExposeEmailId:  emailId,
		Name:           name,
		BotFlag:        botFlag,
		RegisteredDate: registeredDate,
		UpdateDate:     registeredDate,
	}
}

func NewUser(
	persistKey uint,
	identifier text.Identifier,
	name text.Name,
	emailId text.Email,
	botFlag bool,
	registeredDate time.Time,
	updateDate time.Time,
) User {
	return User{
		PersistKey: persistKey,
		UnsavedUser: UnsavedUser{
			Identifier:     identifier,
			ExposeEmailId:  emailId,
			Name:           name,
			BotFlag:        botFlag,
			RegisteredDate: registeredDate,
			UpdateDate:     updateDate,
		},
	}
}
