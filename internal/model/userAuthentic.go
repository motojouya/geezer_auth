package model

import (
	"time"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
)

type CompanyRole struct {
	Company Company
	Roles   []Role
}

type UserAuthentic struc {
	User
	CompanyRole *CompanyRole
	Email       *pkg.Email
}

func NewCompanyRole(company Company, roles []Role) CompanyRole {
	return CompanyRole{
		Company: company,
		Roles:   roles,
	}
}

func NewUserAuthentic(userId uint, exposeId pkg.ExposeId, name pkg.Name, emailId pkg.Email, email *pkg.Email, botFlag bool, registeredDate time.Time, updateDate time.Time, companyRole *CompanyRole) User {
	return UserAuthentic{
		User: NewUser(userId, exposeId, name, emailId, botFlag, registeredDate, updateDate),
		CompanyRole:    companyRole,
		Email:          email,
	}
}

/*
 * CompanyやRoleはpkgをembedして依存関係が明確だが、Userの場合はもうちょいややこしいのとhandlingのtop levelになるのでで変換メソッドがある
 * internal.modelがpkg.modelに依存する形なので、internal.modelに変換関数をもたせる形
 */
(user *UserAuthentic) func ToJwtUser() *pkg.User {
	var companyRole *pkg.CompanyRole = nil
	if user.CompanyRole != nil {
		var company = pkg.NewCompany(user.CompanyRole.Company.ExposeId, user.CompanyRole.Company.Name)

		var roles = make([]*pkg.Role, len(sourceRoles))
		for i, source := range user.CompanyRole.Roles {
			roles[i] = pkg.NewRole(source.Label, source.Name)
		}

		companyRole = pkg.NewCompanyRole(company, roles)
	}

	return pkg.NewUser(
		user.ExposeId,
		user.ExposeEmailId,
		user.Email,
		user.Name,
		user.BotFlag,
		companyRole,
		user.UpdateDate,
	)
}
