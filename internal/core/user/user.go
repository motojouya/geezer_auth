package user

import (
	"time"
	text "github.com/motojouya/geezer_auth/pkg/core/text"
)

const UserExposeIdPrefix = "US-"

type UnsavedUser struct {
	ExposeId       text.ExposeId
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

func CreateUserExposeId(random string) (text.ExposeId, error) {
	return text.CreateExposeId(UserExposeIdPrefix, random)
}

func CreateUser(
	exposeId text.ExposeId,
	emailId text.Email,
	name text.Name,
	botFlag bool,
	registeredDate time.Time
) UnsavedUser {
	return UnsavedUser{
		ExposeId:       exposeId,
		ExposeEmailId:  emailId,
		Name:           name,
		BotFlag:        botFlag,
		RegisteredDate: registeredDate,
		UpdateDate:     registeredDate,
	}
}

func NewUser(
	persistKey uint
	exposeId text.ExposeId,
	name text.Name,
	emailId text.Email,
	botFlag bool,
	registeredDate time.Time,
	updateDate time.Time
) User {
	return User{
		PersistKey:  persistKey,
		UnsavedUser: UnsavedUser{
			ExposeId:       exposeId,
			ExposeEmailId:  emailId,
			Name:           name,
			BotFlag:        botFlag,
			RegisteredDate: registeredDate,
			UpdateDate:     updateDate,
		}
	}
}
