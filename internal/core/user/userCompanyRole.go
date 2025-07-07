package user

import (
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/core/essence"
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
	expireDate *time.Time,
) *UserCompanyRole {
	return &UserCompanyRole{
		PersistKey: persistKey,
		UnsavedUserCompanyRole: UnsavedUserCompanyRole{
			User:         user,
			Company:      company,
			Role:         role,
			RegisterDate: registerDate,
			ExpireDate:   expireDate,
		},
	}
}

func IsUserUCR(user User) func(userCompanyRole *UserCompanyRole) bool {
	return func(userCompanyRole *UserCompanyRole) bool {
		return userCompanyRole.User.Identifier == user.Identifier
	}
}

func SameCompanyUCR(left *UserCompanyRole, right *UserCompanyRole) bool {
	return left.Company.Identifier == right.Company.Identifier
}

func GetRoleUCR(userCompanyRole *UserCompanyRole) role.Role {
	return userCompanyRole.Role
}

func ListToCompanyRole(user User, userCompanyRoles []*UserCompanyRole) (*CompanyRole, error) {

	var allSameUser = essence.Every(userCompanyRoles, IsUserUCR(user))
	if !allSameUser {
		return &CompanyRole{}, essence.NewInvalidArgumentError("UserCompanyRole.User", string(user.Identifier), "UserCompanyRole.User does not match the User")
	}

	var grouped = essence.Group(userCompanyRoles, SameCompanyUCR)
	if len(grouped) > 1 {
		return &CompanyRole{}, essence.NewInvalidArgumentError("UserCompanyRole.Company", "", "UserCompanyRole.Company must be unique for a User")
	}

	var company = userCompanyRoles[0].Company
	var roles = essence.Map(userCompanyRoles, GetRoleUCR)

	return NewCompanyRole(company, roles), nil
}
