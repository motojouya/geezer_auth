package user

import (
	"github.com/doug-martin/goqu/v9"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/db/sql"
	"github.com/motojouya/geezer_auth/internal/db/transfer/company"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"time"
)

type UserCompanyRole struct {
	PersistKey        uint       `db:"persist_key"`
	UserPersistKey    uint       `db:"user_persist_key"`
	CompanyPersistKey uint       `db:"company_persist_key"`
	RoleLabel         string     `db:"role_label"`
	RegisterDate      time.Time  `db:"register_date"`
	ExpireDate        *time.Time `db:"expire_date"`
}

type UserCompanyRoleFull struct {
	UserCompanyRole
	UserIdentifier        string    `db:"user_identifier"`
	UserExposeEmailId     string    `db:"user_email_identifier"`
	UserName              string    `db:"user_name"`
	UserBotFlag           bool      `db:"user_bot_flag"`
	UserRegisteredDate    time.Time `db:"user_register_date"`
	UserUpdateDate        time.Time `db:"user_update_date"`
	CompanyIdentifier     string    `db:"company_identifier"`
	CompanyName           string    `db:"company_name"`
	CompanyRegisteredDate time.Time `db:"company_register_date"`
	RoleName              string    `db:"role_name"`
	RoleDescription       string    `db:"role_description"`
	RoleRegisteredDate    time.Time `db:"role_register_date"`
}

// var SelectUserCompanyRole = sql.Dialect.From("user_company_role").As("ucr").Select(
// 	goqu.C("ucr.persist_key").As("persist_key"),
// 	goqu.C("ucr.user_persist_key").As("user_persist_key"),
// 	goqu.C("ucr.company_persist_key").As("company_persist_key"),
// 	goqu.C("ucr.role_label").As("role_label"),
// 	goqu.C("ucr.register_date").As("register_date"),
// 	goqu.C("ucr.expire_date").As("expire_date"),
// )

var SelectUserCompanyRole = sql.Dialect.From("user_company_role").As("ucr").InnerJoin(
	goqu.T("user").As("u"),
	goqu.On(goqu.Ex{"uat.user_persist_key": goqu.I("u.persist_key")}),
).InnerJoin(
	goqu.T("company").As("c"),
	goqu.On(goqu.Ex{"ci.company_persist_key": goqu.I("c.persist_key")}),
).InnerJoin(
	goqu.T("role").As("r"),
	goqu.On(goqu.Ex{"ci.role_label": goqu.I("r.label")}),
).Select(
	goqu.C("ucr.persist_key").As("persist_key"),
	goqu.C("ucr.user_persist_key").As("user_persist_key"),
	goqu.C("u.identifier").As("user_identifier"),
	goqu.C("u.email_identifier").As("user_email_identifier"),
	goqu.C("u.name").As("user_name"),
	goqu.C("u.bot_flag").As("user_bot_flag"),
	goqu.C("u.register_date").As("user_register_date"),
	goqu.C("u.update_date").As("user_update_date"),
	goqu.C("ucr.company_persist_key").As("company_persist_key"),
	goqu.C("c.identifier").As("company_identifier"),
	goqu.C("c.name").As("company_name"),
	goqu.C("c.register_date").As("company_register_date"),
	goqu.C("ucr.role_label").As("role_label"),
	goqu.C("r.name").As("role_name"),
	goqu.C("r.description").As("role_description"),
	goqu.C("r.register_date").As("role_register_date"),
	goqu.C("ucr.register_date").As("register_date"),
	goqu.C("ucr.expire_date").As("expire_date"),
)

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
