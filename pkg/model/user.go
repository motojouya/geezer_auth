package model

import (
	"time"
)

type User struct {
	ExposeId    ExposeId
	EmailId     Email
	Email       *Email
	Name        Name
	BotFlag     bool
	CompanyRole *CompanyRole
	UpdateDate  time.Time
}

func NewUser(exposeId ExposeId, emailId Email, email *Email, name Name, botFlag bool, companyRole *CompanyRole, updateDate time.Time) *User {
	return &User{
		ExposeId: exposeId,
		EmailId: emailId,
		Email: email,
		Name: name,
		BotFlag: botFlag,
		CompanyRole: companyRole,
		UpdateDate: updateDate,
	}
}
