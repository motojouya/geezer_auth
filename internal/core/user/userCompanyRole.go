package user

import (
	"time"
	"github.com/google/uuid"
	"github.com/motojouya/geezer_auth/internal/model/company"
	"github.com/motojouya/geezer_auth/internal/model/role"
)

type UnsavedUserCompanyRole struct {
	User         User
	Company      company.Company
	Role         role.Role
	RegisterDate time.Time
	ExpireDate   *time.Time
}

type UserCompanyRole struct {
	UserCompanyRoleID  uint
	UnsavedUserCompanyRole
}

func CreateUserCompanyRole(user User, company company.Company, role role.Role, registerDate time.Time) UnsavedUserCompanyRole {
	return UnsavedUserEmail{
		User:         user,
		Company:      company,
		Role:         role,
		RegisterDate: registerDate,
		ExpireDate:   nil,
	}
}

func NewUserEmail(userCompanyRoleID uint, user User, company company.Company, role role.Role, registerDate time.Time, expireDate *time.Time) UserEmail {
	return UserEmail{
		UserCompanyRoleID: userCompanyRoleID,
		UnsavedUserEmail{
			User:         user,
			Company:      company,
			Role:         role,
			RegisterDate: registerDate,
			ExpireDate:   expireDate,
		}
	}
}
