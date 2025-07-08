package user

import (
	core "github.com/motojouya/geezer_auth/internal/core/user"
	text "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
	"github.com/motojouya/geezer_auth/internal/db"
	"github.com/doug-martin/goqu/v9"
)

type UserAuthentic struct {
	UserPersistKey     uint
	UserIdentifier     string
	UserExposeEmailId  string
	UserName           string
	UserBotFlag        bool
	UserRegisteredDate time.Time
	UserUpdateDate     time.Time
	Email              *string
	UserCompanyRole    []*UserCompanyRoleFull
}

// TODO user join user_email
var SelectFullUserCompanyRole = db.Dialect.From("user").As("u").LeftOuterJoin(
	goqu.T("user_email").As("ue"),
	goqu.On(goqu.Ex{"uat.user_persist_key": goqu.I("ue.persist_key")}),
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

func RelateUserCompanyRole(ua *UserAuthentic, ucr *UserCompanyRoleFull) (*UserAuthentic, bool) {
	if ua.UserPersistKey == ucr.UserPersistKey {
		ua.UserCompanyRole = append(ua.UserCompanyRole, ucr)
		return ua, true
	} else {
		return ua, false
	}
}

func (ua UserAuthentic) ToCoreUserAuthentic() (*core.UserAuthentic, error) {
	var user, userErr = (User{
		PersistKey:     ua.UserPersistKey,
		Identifier:     ua.UserIdentifier,
		ExposeEmailId:  ua.UserExposeEmailId,
		Name:           ua.UserName,
		BotFlag:        ua.UserBotFlag,
		RegisteredDate: ua.UserRegisteredDate,
		UpdateDate:     ua.UserUpdateDate,
	}).ToCoreUser()
	if userErr != nil {
		return &core.UserAuthentic{}, userErr
	}

	var email *text.Email = nil
	if ua.Email != nil {
		var emailResult, emailErr = text.NewEmail(*ua.Email)
		if emailErr != nil {
			return &core.UserAuthentic{}, emailErr
		}
		email = &emailResult
	}

	var coreUserCompanyRoles = make([]*core.UserCompanyRole, 0, len(ua.UserCompanyRole))
	for _, ucr := range ua.UserCompanyRole {
		var coreUserCompanyRole, companyRoleErr = ucr.ToCoreUserCompanyRole()
		if companyRoleErr != nil {
			return &core.UserAuthentic{}, companyRoleErr
		}
		coreUserCompanyRoles = append(coreUserCompanyRoles, coreUserCompanyRole)
	}

	var companyRole, companyRoleErr = core.ListToCompanyRole(user, coreUserCompanyRoles)
	if companyRoleErr != nil {
		return &core.UserAuthentic{}, companyRoleErr
	}

	return core.NewUserAuthentic(
		user,
		companyRole,
		email,
	), nil
}
