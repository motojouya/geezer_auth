package common

import (
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	"time"
)

type CompanyRole struct {
	Company Company
	Roles   []Role
}

type User struct {
	Identifier      string
	IdentifierEmail string
	Name            string
	BotFlag         bool
	UpdateDate      time.Time
	CompanyRole     *CompanyRole
	Email           *string
}

func FromShelterUser(u *shelter.UserAuthentic) *User {
	var companyRole *CompanyRole = nil
	if u.CompanyRole != nil {
		companyRole = &CompanyRole{
			Company: FromShelterCompany(u.CompanyRole.Company),
			Roles:   essence.Map(u.CompanyRole.Roles, FromShelterRole),
		}
	}

	var email *string = nil
	if u.Email != nil {
		var emailStr = string(*u.Email)
		email = &emailStr
	}

	return &User{
		Identifier:      string(u.Identifier),
		IdentifierEmail: string(u.ExposeEmailId),
		Name:            string(u.Name),
		BotFlag:         u.BotFlag,
		UpdateDate:      u.UpdateDate,
		CompanyRole:     companyRole,
		Email:           email,
	}
}
