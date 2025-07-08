package user

import (
	"github.com/doug-martin/goqu/v9"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/db"
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type UserAccessToken struct {
	PersistKey       uint      `db:"persist_key"`
	UserPersistKey   uint      `db:"user_persist_key"`
	AccessToken      string    `db:"access_token"`
	SourceUpdateDate time.Time `db:"source_update_date"`
	RegisterDate     time.Time `db:"register_date"`
	ExpireDate       time.Time `db:"expire_date"`
}

type UserAccessTokenFull struct {
	UserAccessToken
	UserIdentifier     string    `db:"user_identifier"`
	UserExposeEmailId  string    `db:"user_email_identifier"`
	UserName           string    `db:"user_name"`
	UserBotFlag        bool      `db:"user_bot_flag"`
	UserRegisteredDate time.Time `db:"user_register_date"`
	UserUpdateDate     time.Time `db:"user_update_date"`
}

var SelectUserAccessToken = db.Dialect.From("user_access_token").As("uat").Select(
	goqu.C("uat.persist_key").As("persist_key"),
	goqu.C("uat.user_persist_key").As("user_persist_key"),
	goqu.C("uat.access_token").As("access_token"),
	goqu.C("uat.source_update_date").As("source_update_date"),
	goqu.C("uat.register_date").As("register_date"),
	goqu.C("uat.expire_date").As("expire_date"),
)

var SelectFullUserAccessToken = db.Dialect.From("user_access_token").As("uat").InnerJoin(
	goqu.T("user").As("u"),
	goqu.On(goqu.Ex{"uat.user_persist_key": goqu.I("u.persist_key")}),
).Select(
	goqu.C("uat.persist_key").As("persist_key"),
	goqu.C("uat.user_persist_key").As("user_persist_key"),
	goqu.C("u.identifier").As("user_identifier"),
	goqu.C("u.email_identifier").As("user_email_identifier"),
	goqu.C("u.name").As("user_name"),
	goqu.C("u.bot_flag").As("user_bot_flag"),
	goqu.C("u.register_date").As("user_register_date"),
	goqu.C("u.update_date").As("user_update_date"),
	goqu.C("uat.access_token").As("access_token"),
	goqu.C("uat.source_update_date").As("source_update_date"),
	goqu.C("uat.register_date").As("register_date"),
	goqu.C("uat.expire_date").As("expire_date"),
)

func FromCoreUserAccessToken(coreUserAccessToken core.UnsavedUserAccessToken) UserAccessToken {
	return UserAccessToken{
		UserPersistKey:   coreUserAccessToken.User.PersistKey,
		AccessToken:      string(coreUserAccessToken.AccessToken),
		SourceUpdateDate: coreUserAccessToken.SourceUpdateDate,
		RegisterDate:     coreUserAccessToken.RegisterDate,
		ExpireDate:       coreUserAccessToken.ExpireDate,
	}
}

func (ua UserAccessTokenFull) ToCoreUserAccessToken() (core.UserAccessToken, error) {
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
		return core.UserAccessToken{}, userErr
	}

	return core.NewUserAccessToken(
		ua.PersistKey,
		user,
		text.NewJwtToken(ua.AccessToken),
		ua.SourceUpdateDate,
		ua.RegisterDate,
		ua.ExpireDate,
	), nil
}
