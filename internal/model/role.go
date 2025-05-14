package model

import (
	"time"
)

type UnsavedRole struct {
	name        string
	label       string
	description string
}

type Role struct {
	roleId         uint
	registeredDate time.Time
	UnsavedRole
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

func CreateRole(name string, label string, description string) UnsavedRole {
	return UnsavedRole{
		name:        name,
		label:       label,
		description: description,
	}
}

func NewRole(roleId uint, name string, label string, description string, registeredDate time.Time) Role {
	return Role{
		roleId:         roleId,
		registeredDate: registeredDate,
		UnsavedRole:    CreateRole(name, label, description),
	}
}
