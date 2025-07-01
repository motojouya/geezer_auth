package user

import (
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/role"
	utility "github.com/motojouya/geezer_auth/internal/utility"
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

	var allSameUser = utility.Every(userCompanyRoles, IsUserUCR(user))
	if !allSameUser {
		return &CompanyRole{}, utility.NewInvalidArgumentError("UserCompanyRole.User", string(user.Identifier), "UserCompanyRole.User does not match the User")
	}

	var grouped = utility.Group(userCompanyRoles, SameCompanyUCR)
	if len(grouped) > 1 {
		return &CompanyRole{}, utility.NewInvalidArgumentError("UserCompanyRole.Company", "", "UserCompanyRole.Company must be unique for a User")
	}

	var company = userCompanyRoles[0].Company
	var roles = utility.Map(userCompanyRoles, GetRoleUCR)

	return NewCompanyRole(company, roles), nil
}
