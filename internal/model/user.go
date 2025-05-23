package model

import (
	"time"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
)

const UserExposeIdPrefix = "US-"

type UnsavedUser struct {
	ExposeId       pkg.ExposeId
	ExposeEmailId  pkg.Email
	Name           pkg.Name
	BotFlag        bool
	RegisteredDate time.Time
	UpdateDate     time.Time
}

type User struct {
	UserId         uint
	CompanyRole    *CompanyRole
	Email          *pkg.Email
	UnsavedUser
}

func CreateUserExposeId(random string) (pkg.ExposeId, error) {
	return pkg.CreateExposeId(UserExposeIdPrefix, random)
}

func CreateUser(exposeId pkg.ExposeId, emailId pkg.Email, name pkg.Name, botFlag bool, registeredDate time.Time, updateDate time.Time) UnsavedUser {
	return UnsavedUser{
		ExposeId:       exposeId,
		ExposeEmailId:  emailId,
		Name:           name,
		BotFlag:        botFlag,
		RegisteredDate: registeredDate,
		UpdateDate:     updateDate,
	}
}

func NewUser(userId uint, exposeId pkg.ExposeId, name pkg.Name, emailId pkg.Email, email *pkg.Email, botFlag bool, registeredDate time.Time, updateDate time.Time, companyRole *CompanyRole) User {
	return User{
		UserId:         userId,
		CompanyRole:    companyRole,
		Email:          email,
		UnsavedUser:    CreateUser(exposeId, emailId, name, botFlag, registeredDate, updateDate),
	}
}

/*
 * CompanyやRoleはpkgをembedして依存関係が明確だが、Userの場合はもうちょいややこしいのとhandlingのtop levelになるのでで変換メソッドがある
 * internal.modelがpkg.modelに依存する形なので、internal.modelに変換関数をもたせる形
 */
(user *User) func ToJwtUser() *pkg.User {
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
