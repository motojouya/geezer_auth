package model

import (
	"time"
	"github.com/google/uuid"
)

type UnsavedCompanyInvite struct {
	Company        Company
	User           *User
	Role           *Role
	Token          uuid.UUID
	RegisteredDate time.Time
	ExpireDate     *time.Time
}

type CompanyInvite struct {
	CompanyInviteId uint
	UnsavedCompanyInvite
}

func CreateCompanyInvite(company Company, token uuid.UUID, user *User, role *Role, registerDate time.Time) UnsavedCompanyInvite {
	return UnsavedCompanyInvite{
		Company:      company,
		User:         user,
		Role:         role,
		RegisterDate: registerDate,
		ExpireDate:   nil,
	}
}

func NewUserRefreshToken(companyInviteId uint, company Company, token uuid.UUID, user *User, role *Role, registerDate time.Time, expireDate *time.Time) CompanyInvite {
	return CompanyInvite{
		CompanyInviteId: companyInviteId,
		UnsavedCompanyInvite: UnsavedCompanyInvite{
			Company:      company,
			Token:        token,
			User:         user,
			Role:         role,
			RegisterDate: registerDate,
			ExpireDate:   expireDate,
		},
	}
}
