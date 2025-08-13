package user

import (
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"time"
)

const UserIdentifierPrefix = "US-"

type UnsavedUser struct {
	Identifier     pkgText.Identifier
	ExposeEmailId  pkgText.Email
	Name           pkgText.Name
	BotFlag        bool
	RegisteredDate time.Time
	UpdateDate     time.Time
}

type User struct {
	PersistKey uint
	UnsavedUser
}

func CreateUserIdentifier(random string) (pkgText.Identifier, error) {
	return pkgText.CreateIdentifier(UserIdentifierPrefix, random)
}

func CreateUser(
	identifier pkgText.Identifier,
	emailId pkgText.Email,
	name pkgText.Name,
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
	identifier pkgText.Identifier,
	name pkgText.Name,
	emailId pkgText.Email,
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

func (user User) Update(updateDate time.Time) User {
	return User{
		PersistKey: user.PersistKey,
		UnsavedUser: UnsavedUser{
			Identifier:     user.Identifier,
			ExposeEmailId:  user.ExposeEmailId,
			Name:           user.Name,
			BotFlag:        user.BotFlag,
			RegisteredDate: user.RegisteredDate,
			UpdateDate:     updateDate,
		},
	}
}

func (user User) UpdateName(name pkgText.Name, updateDate time.Time) User {
	return User{
		PersistKey: user.PersistKey,
		UnsavedUser: UnsavedUser{
			Identifier:     user.Identifier,
			ExposeEmailId:  user.ExposeEmailId,
			Name:           name,
			BotFlag:        user.BotFlag,
			RegisteredDate: user.RegisteredDate,
			UpdateDate:     updateDate,
		},
	}
}
