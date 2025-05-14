package model

import (
	"time"
)

const UserExposeIdPrefix = "US-"

type UnsavedUser struct {
	ExposeId      string
	ExposeEmailId string
	Name          string
	BotFlag       bool
}

type User struct {
	UserId         uint
	CompanyRole    *CompanyRole
	Email          *string
	RegisteredDate time.Time
	UnsavedUser
}

func CreateUser(exposeId string, emailId string, name string, botFlag bool) UnsavedUser {
	return UnsavedUser{
		ExposeId:      exposeId,
		ExposeEmailId: emailId,
		Name:          name,
		BotFlag:       botFlag,
	}
}

func NewUser(userId uint, exposeId string, name string, emailId string, email *string, botFlag bool, registeredDate time.Time, companyRole *CompanyRole) User {
	return User{
		UserId:         userId,
		CompanyRole:    companyRole,
		Email:          email,
		RegisteredDate: registeredDate,
		UnsavedUser:    CreateUser(exposeId, name, emailId, botFlag),
	}
}
