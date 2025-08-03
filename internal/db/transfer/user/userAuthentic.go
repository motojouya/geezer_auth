package user

import (
	"github.com/doug-martin/goqu/v9"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	text "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"time"
)

type UserAuthentic struct {
	UserPersistKey     uint                  `db:"persist_key"`
	UserIdentifier     string                `db:"identifier"`
	UserExposeEmailId  string                `db:"email_identifier"`
	UserName           string                `db:"name"`
	UserBotFlag        bool                  `db:"bot_flag"`
	UserRegisteredDate time.Time             `db:"register_date"`
	UserUpdateDate     time.Time             `db:"update_date"`
	Email              *string               `db:"email"`
	UserCompanyRole    []UserCompanyRoleFull `db:"-"`
}

var SelectUserAuthentic = utility.Dialect.From(goqu.T("users").As("u")).LeftOuterJoin(
	goqu.T("user_email").As("ue"),
	goqu.On(
		goqu.I("u.persist_key").Eq(goqu.I("ue.user_persist_key")),
		goqu.I("ue.verify_date").IsNotNull(),
		goqu.I("ue.expire_date").IsNull(),
	),
).Select(
	goqu.I("u.persist_key").As("persist_key"),
	goqu.I("u.identifier").As("identifier"),
	goqu.I("u.email_identifier").As("email_identifier"),
	goqu.I("u.name").As("name"),
	goqu.I("u.bot_flag").As("bot_flag"),
	goqu.I("u.register_date").As("register_date"),
	goqu.I("u.update_date").As("update_date"),
	goqu.I("ue.email").As("email"),
)

func UserIdentifierUserAuthentic(userAuthentic *UserAuthentic) string {
	return userAuthentic.UserIdentifier
}

func RelateUserCompanyRole(ua *UserAuthentic, ucr UserCompanyRoleFull) (*UserAuthentic, bool) {
	if ua.UserPersistKey == ucr.UserPersistKey {
		ua.UserCompanyRole = append(ua.UserCompanyRole, ucr)
		return ua, true
	} else {
		return ua, false
	}
}

func (ua UserAuthentic) ToCoreUserAuthentic() (*shelter.UserAuthentic, error) {
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
		return &shelter.UserAuthentic{}, userErr
	}

	var email *text.Email = nil
	if ua.Email != nil {
		var emailResult, emailErr = text.NewEmail(*ua.Email)
		if emailErr != nil {
			return &shelter.UserAuthentic{}, emailErr
		}
		email = &emailResult
	}

	var shelterUserCompanyRoles = make([]shelter.UserCompanyRole, 0, len(ua.UserCompanyRole))
	for _, ucr := range ua.UserCompanyRole {
		var shelterUserCompanyRole, companyRoleErr = ucr.ToCoreUserCompanyRole()
		if companyRoleErr != nil {
			return &shelter.UserAuthentic{}, companyRoleErr
		}
		shelterUserCompanyRoles = append(shelterUserCompanyRoles, *shelterUserCompanyRole)
	}

	var companyRole, companyRoleErr = shelter.ListToCompanyRole(user, shelterUserCompanyRoles)
	if companyRoleErr != nil {
		return &shelter.UserAuthentic{}, companyRoleErr
	}

	return shelter.NewUserAuthentic(
		user,
		companyRole,
		email,
	), nil
}
