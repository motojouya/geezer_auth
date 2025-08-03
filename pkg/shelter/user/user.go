package user

import (
	"github.com/motojouya/geezer_auth/pkg/shelter/text"
	"time"
)

type User struct {
	Identifier  text.Identifier
	EmailId     text.Email
	Email       *text.Email
	Name        text.Name
	BotFlag     bool
	CompanyRole *CompanyRole
	UpdateDate  time.Time
}

func NewUser(
	identifier text.Identifier,
	emailId text.Email,
	email *text.Email,
	name text.Name,
	botFlag bool,
	companyRole *CompanyRole,
	updateDate time.Time,
) *User {
	return &User{
		Identifier:  identifier,
		EmailId:     emailId,
		Email:       email,
		Name:        name,
		BotFlag:     botFlag,
		CompanyRole: companyRole,
		UpdateDate:  updateDate,
	}
}
