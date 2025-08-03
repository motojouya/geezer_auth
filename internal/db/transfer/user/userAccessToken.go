package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"github.com/motojouya/geezer_auth/pkg/shelter/text"
	"time"
)

type UserAccessToken struct {
	PersistKey       uint      `db:"persist_key,primarykey,autoincrement"`
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

func AddUserAccessTokenTable(dbMap *gorp.DbMap) {
	dbMap.AddTableWithName(UserAccessToken{}, "user_access_token").SetKeys(true, "PersistKey")
}

// var SelectUserAccessToken = utility.Dialect.From("user_access_token").As("uat").Select(
// 	goqu.C("uat.persist_key").As("persist_key"),
// 	goqu.C("uat.user_persist_key").As("user_persist_key"),
// 	goqu.C("uat.access_token").As("access_token"),
// 	goqu.C("uat.source_update_date").As("source_update_date"),
// 	goqu.C("uat.register_date").As("register_date"),
// 	goqu.C("uat.expire_date").As("expire_date"),
// )

var SelectUserAccessToken = utility.Dialect.From(goqu.T("user_access_token").As("uat")).InnerJoin(
	goqu.T("users").As("u"),
	goqu.On(goqu.Ex{"uat.user_persist_key": goqu.I("u.persist_key")}),
).Select(
	goqu.I("uat.persist_key").As("persist_key"),
	goqu.I("uat.user_persist_key").As("user_persist_key"),
	goqu.I("u.identifier").As("user_identifier"),
	goqu.I("u.email_identifier").As("user_email_identifier"),
	goqu.I("u.name").As("user_name"),
	goqu.I("u.bot_flag").As("user_bot_flag"),
	goqu.I("u.register_date").As("user_register_date"),
	goqu.I("u.update_date").As("user_update_date"),
	goqu.I("uat.access_token").As("access_token"),
	goqu.I("uat.source_update_date").As("source_update_date"),
	goqu.I("uat.register_date").As("register_date"),
	goqu.I("uat.expire_date").As("expire_date"),
)

func FromCoreUserAccessToken(shelterUserAccessToken shelter.UnsavedUserAccessToken) UserAccessToken {
	return UserAccessToken{
		UserPersistKey:   shelterUserAccessToken.User.PersistKey,
		AccessToken:      string(shelterUserAccessToken.AccessToken),
		SourceUpdateDate: shelterUserAccessToken.SourceUpdateDate,
		RegisterDate:     shelterUserAccessToken.RegisterDate,
		ExpireDate:       shelterUserAccessToken.ExpireDate,
	}
}

func (ua UserAccessTokenFull) ToCoreUserAccessToken() (shelter.UserAccessToken, error) {
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
		return shelter.UserAccessToken{}, userErr
	}

	return shelter.NewUserAccessToken(
		ua.PersistKey,
		user,
		text.NewJwtToken(ua.AccessToken),
		ua.SourceUpdateDate,
		ua.RegisterDate,
		ua.ExpireDate,
	), nil
}

// testdata投入時に楽するためのもの。アプリケーションからは利用を想定しない。
func NewUserAccessToken(persistKey uint, userPersistKey uint, accessToken string, sourceUpdateDate time.Time, registerDate time.Time, expireDate time.Time) UserAccessToken {
	return UserAccessToken{
		PersistKey:       persistKey,
		UserPersistKey:   userPersistKey,
		AccessToken:      accessToken,
		SourceUpdateDate: sourceUpdateDate,
		RegisterDate:     registerDate,
		ExpireDate:       expireDate,
	}
}
