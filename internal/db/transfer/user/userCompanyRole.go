package user

import (
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/db/transfer/company"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"time"
)

type UserCompanyRole struct {
	PersistKey        uint
	UserPersistKey    uint
	CompanyPersistKey uint
	RoleLabel         string
	RegisterDate      time.Time
	ExpireDate        *time.Time
}

type UserCompanyRoleFull struct {
	UserCompanyRole
	UserIdentifier        string
	UserExposeEmailId     string
	UserName              string
	UserBotFlag           bool
	UserRegisteredDate    time.Time
	UserUpdateDate        time.Time
	CompanyIdentifier     string
	CompanyName           string
	CompanyRegisteredDate time.Time
	RoleName              string
	RoleDescription       string
	RoleRegisteredDate    time.Time
}

func FromCoreUserCompanyRole(coreUserCompanyRole *core.UnsavedUserCompanyRole) *UserCompanyRole {
	return &UserCompanyRole{
		UserPersistKey:    coreUserCompanyRole.User.PersistKey,
		CompanyPersistKey: coreUserCompanyRole.Company.PersistKey,
		RoleLabel:         string(coreUserCompanyRole.Role.Label),
		RegisterDate:      coreUserCompanyRole.RegisterDate,
		ExpireDate:        coreUserCompanyRole.ExpireDate,
	}
}

func (u UserCompanyRoleFull) ToCoreUserCompanyRole() (*core.UserCompanyRole, error) {
	var user, userErr = (User{
		PersistKey:     u.UserPersistKey,
		Identifier:     u.UserIdentifier,
		ExposeEmailId:  u.UserExposeEmailId,
		Name:           u.UserName,
		BotFlag:        u.UserBotFlag,
		RegisteredDate: u.UserRegisteredDate,
		UpdateDate:     u.UserUpdateDate,
	}).ToCoreUser()
	if userErr != nil {
		return &core.UserCompanyRole{}, userErr
	}

	var companyValue, companyErr = (company.Company{
		PersistKey:     u.CompanyPersistKey,
		Identifier:     u.CompanyIdentifier,
		Name:           u.CompanyName,
		RegisteredDate: u.CompanyRegisteredDate,
	}).ToCoreCompany()
	if companyErr != nil {
		return &core.UserCompanyRole{}, companyErr
	}

	var roleValue, roleErr = (role.Role{
		Label:          u.RoleLabel,
		Name:           u.RoleName,
		Description:    u.RoleDescription,
		RegisteredDate: u.RoleRegisteredDate,
	}).ToCoreRole()
	if roleErr != nil {
		return &core.UserCompanyRole{}, roleErr
	}

	return core.NewUserCompanyRole(
		u.PersistKey,
		user,
		companyValue,
		roleValue,
		u.UserCompanyRole.RegisterDate,
		u.UserCompanyRole.ExpireDate,
	), nil
}
