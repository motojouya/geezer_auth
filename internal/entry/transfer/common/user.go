package common

import (
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	"time"
)

type CompanyRole struct {
	Company Company `json:"company"`
	Roles   []Role  `json:"roles"`
}

type User struct {
	Identifier      string       `json:"identifier"`
	IdentifierEmail string       `json:"identifier_email"`
	Name            string       `json:"name"`
	BotFlag         bool         `json:"bot_flag"`
	UpdateDate      time.Time    `json:"update_date"`
	CompanyRole     *CompanyRole `json:"company_role"`
	Email           *string      `json:"email"`
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
