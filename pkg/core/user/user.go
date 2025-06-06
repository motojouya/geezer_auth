package user

import (
	"time"
)

type User struct {
	Identifier    Identifier
	EmailId     Email
	Email       *Email
	Name        Name
	BotFlag     bool
	CompanyRole *CompanyRole
	UpdateDate  time.Time
}

func NewUser(
	identifier    Identifier,
	emailId     Email,
	email       *Email,
	name        Name,
	botFlag     bool,
	companyRole *CompanyRole,
	updateDate  time.Time
) *User {
	return &User{
		Identifier: identifier,
		EmailId: emailId,
		Email: email,
		Name: name,
		BotFlag: botFlag,
		CompanyRole: companyRole,
		UpdateDate: updateDate,
	}
}
