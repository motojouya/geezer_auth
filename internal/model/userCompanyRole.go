package model

import (
	"time"
	"github.com/google/uuid"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
)

type UnsavedUserCompanyRole struct {
	User         User
	Company      Company
	Role         Role
	RegisterDate time.Time
	ExpireDate   *time.Time
}

type UserCompanyRole struct {
	UserCompanyRoleID  uint
	UnsavedUserCompanyRole
}

func CreateUserCompanyRole(user User, company Company, role Role, registerDate time.Time) UnsavedUserCompanyRole {
	return UnsavedUserEmail{
		User:         user,
		Company:      company,
		Role:         role,
		RegisterDate: registerDate,
		ExpireDate:   nil,
	}
}

func NewUserEmail(userCompanyRoleID uint, user User, company Company, role Role, registerDate time.Time, expireDate *time.Time) UserEmail {
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
