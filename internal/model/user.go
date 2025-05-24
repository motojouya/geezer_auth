package model

import (
	"time"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
)

const UserExposeIdPrefix = "US-"

type UnsavedUser struct {
	ExposeId       pkg.ExposeId
	ExposeEmailId  pkg.Email
	Name           pkg.Name
	BotFlag        bool
	RegisteredDate time.Time
	UpdateDate     time.Time
}

type User struct {
	UserId uint
	UnsavedUser
}

func CreateUserExposeId(random string) (pkg.ExposeId, error) {
	return pkg.CreateExposeId(UserExposeIdPrefix, random)
}

func CreateUser(exposeId pkg.ExposeId, emailId pkg.Email, name pkg.Name, botFlag bool, registeredDate time.Time) UnsavedUser {
	return UnsavedUser{
		ExposeId:       exposeId,
		ExposeEmailId:  emailId,
		Name:           name,
		BotFlag:        botFlag,
		RegisteredDate: registeredDate,
		UpdateDate:     registeredDate,
	}
}

func NewUser(userId uint, exposeId pkg.ExposeId, name pkg.Name, emailId pkg.Email, botFlag bool, registeredDate time.Time, updateDate time.Time) User {
	return User{
		UserId:         userId,
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

(user *User) func SetName(name pkg.Name) {
	user.Name = name
}
