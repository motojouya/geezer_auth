package model

import (
	"time"
)

const UserExposeIdPrefix = "US-"

type UnsavedUser struct {
	ExposeId      ExposeId
	ExposeEmailId Email
	Name          Name
	BotFlag       bool
}

type User struct {
	UserId         uint
	CompanyRole    *CompanyRole
	Email          *Email
	RegisteredDate time.Time
	UpdateDate     time.Time
	UnsavedUser
}

func NewUserExposeId(random string) (ExposeId, error) {
	return NewExposeId(UserExposeIdPrefix, random)
}

func CreateUser(exposeId ExposeId, emailId Email, name Name, botFlag bool) UnsavedUser {
	return UnsavedUser{
		ExposeId:      exposeId,
		ExposeEmailId: emailId,
		Name:          name,
		BotFlag:       botFlag,
	}
}

func NewUser(userId uint, exposeId ExposeId, name Name, emailId Email, email *Email, botFlag bool, registeredDate time.Time, updateDate time.Time, companyRole *CompanyRole) User {
	return User{
		UserId:         userId,
		CompanyRole:    companyRole,
		Email:          email,
		RegisteredDate: registeredDate,
		UpdateDate:     updateDate,
		UnsavedUser:    CreateUser(exposeId, emailId, name, botFlag),
	}
}
