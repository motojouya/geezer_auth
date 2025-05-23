package model

import (
	"time"
	"github.com/google/uuid"
)

type UnsavedCompanyInvite struct {
	Company        Company
	Token          uuid.UUID
	Role           Role
	User           *User
	RegisteredDate time.Time
	ExpireDate     time.Time
}

type CompanyInvite struct {
	CompanyInviteId uint
	UnsavedCompanyInvite
}

const TokenValidityPeriodHours = 50

func CreateCompanyInvite(company Company, token uuid.UUID, role Role, user *User, registerDate time.Time) UnsavedCompanyInvite {
	var expireDate = registerDate.Add(TokenValidityPeriodHours * time.Hour)

	return UnsavedCompanyInvite{
		Company:      company,
		Token:        token,
		Role:         role,
		User:         user,
		RegisterDate: registerDate,
		ExpireDate:   expireDate,
	}
}

func NewUserRefreshToken(companyInviteId uint, company Company, token uuid.UUID, role Role, user *User, registerDate time.Time, expireDate time.Time) CompanyInvite {
	return CompanyInvite{
		CompanyInviteId: companyInviteId,
		UnsavedCompanyInvite: UnsavedCompanyInvite{
			Company:      company,
			Token:        token,
			Role:         role,
			User:         user,
			RegisterDate: registerDate,
			ExpireDate:   expireDate,
		},
	}
}
