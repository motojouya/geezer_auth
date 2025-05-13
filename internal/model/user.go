package model

import (
	"time"
)

const UserExposeIdPrefix = "US-"

// company CO-

type UnsavedUser struct {
	exposeId      string
	exposeEmailId string
	name          string
	botFlag       bool
}

type User struct {
	userId         uint
	companyRole    *CompanyRole
	email          *string
	registeredDate time.Time
	UnsavedUser
}

func CreateUser(exposeId string, emailId string, name string, botFlag bool) UnsavedUser {
	return UnsavedUser{
		exposeId:      exposeId,
		exposeEmailId: emailId,
		name:          name,
		botFlag:       botFlag,
	}
}

func NewUser(userId uint, exposeId string, name string, emailId string, email *string, botFlag bool, registeredDate time.Time, companyRole *CompanyRole) User {
	return User{
		userId:         userId,
		companyRole:    companyRole,
		email:          email,
		registeredDate: registeredDate,
		UnsavedUser:    CreateUser(exposeId, name, emailId, botFlag),
	}
}
