package user

import (
	"time"
	text "github.com/motojouya/geezer_auth/pkg/model/text"
	pkg "github.com/motojouya/geezer_auth/pkg/model/user"
	"github.com/motojouya/geezer_auth/internal/model/company"
	"github.com/motojouya/geezer_auth/internal/model/role"
)

type CompanyRole struct {
	Company company.Company
	Roles   []role.Role
}

type UserAuthentic struc {
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
	userId uint,
	exposeId text.ExposeId,
	name text.Name,
	emailId text.Email,
	email *text.Email,
	botFlag bool,
	registeredDate time.Time,
	updateDate time.Time,
	companyRole *CompanyRole
) *UserAuthentic {
	return &UserAuthentic{
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
