package user

import (
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/role"
	"time"
)

type UnsavedUserCompanyRole struct {
	User         User
	Company      company.Company
	Role         role.Role
	RegisterDate time.Time
	ExpireDate   *time.Time
}

type UserCompanyRole struct {
	PersistKey uint
	UnsavedUserCompanyRole
}

func CreateUserCompanyRole(
	user User,
	company company.Company,
	role role.Role,
	registerDate time.Time,
) *UnsavedUserCompanyRole {
	return &UnsavedUserCompanyRole{
		User:         user,
		Company:      company,
		Role:         role,
		RegisterDate: registerDate,
		ExpireDate:   nil,
	}
}

func NewUserCompanyRole(
	persistKey uint,
	user User,
	company company.Company,
	role role.Role,
	registerDate time.Time,
	expireDate time.Time,
) *UserCompanyRole {
	return &UserCompanyRole{
		PersistKey:             persistKey,
		UnsavedUserCompanyRole: UnsavedUserCompanyRole{
			User:         user,
			Company:      company,
			Role:         role,
			RegisterDate: registerDate,
			ExpireDate:   &expireDate,
		},
	}
}
