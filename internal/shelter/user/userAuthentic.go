package user

import (
	"github.com/motojouya/geezer_auth/internal/shelter/company"
	"github.com/motojouya/geezer_auth/internal/shelter/role"
	text "github.com/motojouya/geezer_auth/pkg/shelter/text"
	pkg "github.com/motojouya/geezer_auth/pkg/shelter/user"
)

type CompanyRole struct {
	Company company.Company
	Roles   []role.Role
}

type UserAuthentic struct {
	User
	CompanyRole *CompanyRole
	Email       *text.Email
}

func NewCompanyRole(company company.Company, roles []role.Role) *CompanyRole {
	return &CompanyRole{
		Company: company,
		Roles:   roles,
	}
}

func NewUserAuthentic(
	user User,
	companyRole *CompanyRole,
	email *text.Email,
) *UserAuthentic {
	return &UserAuthentic{
		User:        user,
		CompanyRole: companyRole,
		Email:       email,
	}
}

/*
 * CompanyやRoleはpkgをembedして依存関係が明確だが、Userの場合はもうちょいややこしいのとhandlingのtop levelになるのでで変換メソッドがある
 * internal.modelがpkg.modelに依存する形なので、internal.modelに変換関数をもたせる形
 */
func (user *UserAuthentic) ToJwtUser() *pkg.User {
	var companyRole *pkg.CompanyRole = nil
	if user.CompanyRole != nil {
		var company = pkg.NewCompany(user.CompanyRole.Company.Identifier, user.CompanyRole.Company.Name)

		var roles = make([]pkg.Role, len(user.CompanyRole.Roles))
		for i, source := range user.CompanyRole.Roles {
			roles[i] = pkg.NewRole(source.Label, source.Name)
		}

		companyRole = pkg.NewCompanyRole(company, roles)
	}

	return pkg.NewUser(
		user.Identifier,
		user.ExposeEmailId,
		user.Email,
		user.Name,
		user.BotFlag,
		companyRole,
		user.UpdateDate,
	)
}

func (user *UserAuthentic) GetUser() User {
	return user.User
}
