package user

import (
	"github.com/motojouya/geezer_auth/internal/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"time"
	"github.com/motojouya/geezer_auth/internal/db"
	"github.com/doug-martin/goqu/v9"
)

type UserRefreshToken struct {
	PersistKey     uint
	UserPersistKey uint
	RefreshToken   string
	RegisterDate   time.Time
	ExpireDate     time.Time
}

type UserRefreshTokenFull struct {
	UserRefreshToken
	UserIdentifier     string
	UserExposeEmailId  string
	UserName           string
	UserBotFlag        bool
	UserRegisteredDate time.Time
	UserUpdateDate     time.Time
}

var SelectUserRefreshToken = db.Dialect.From("user_refresh_token").As("urt").Select(
	goqu.C("urt.persist_key").As("persist_key"),
	goqu.C("urt.user_persist_key").As("user_persist_key"),
	goqu.C("urt.refresh_token").As("refresh_token"),
	goqu.C("urt.register_date").As("register_date"),
	goqu.C("urt.expire_date").As("expire_date"),
)

var SelectFullUserRefreshToken = db.Dialect.From("user_refresh_token").As("urt").InnerJoin(
	goqu.T("user").As("u"),
	goqu.On(goqu.Ex{"uat.user_persist_key": goqu.I("u.persist_key")}),
).Select(
	goqu.C("urt.persist_key").As("persist_key"),
	goqu.C("urt.user_persist_key").As("user_persist_key"),
	goqu.C("u.identifier").As("user_identifier"),
	goqu.C("u.email_identifier").As("user_email_identifier"),
	goqu.C("u.name").As("user_name"),
	goqu.C("u.bot_flag").As("user_bot_flag"),
	goqu.C("u.register_date").As("user_register_date"),
	goqu.C("u.update_date").As("user_update_date"),
	goqu.C("urt.refresh_token").As("refresh_token"),
	goqu.C("urt.register_date").As("register_date"),
	goqu.C("urt.expire_date").As("expire_date"),
)

func FromCoreUserRefreshToken(coreUserRefreshToken core.UnsavedUserRefreshToken) UserRefreshToken {
	return UserRefreshToken{
		UserPersistKey: coreUserRefreshToken.User.PersistKey,
		RefreshToken:   string(coreUserRefreshToken.RefreshToken),
		RegisterDate:   coreUserRefreshToken.RegisterDate,
		ExpireDate:     coreUserRefreshToken.ExpireDate,
	}
}

func (u UserRefreshTokenFull) ToCoreUserRefreshToken() (core.UserRefreshToken, error) {
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
		return core.UserRefreshToken{}, userErr
	}

	var token, tokenErr = text.NewToken(u.RefreshToken)
	if tokenErr != nil {
		return core.UserRefreshToken{}, tokenErr
	}

	return core.NewUserRefreshToken(
		u.PersistKey,
		user,
		token,
		u.RegisterDate,
		u.ExpireDate,
	), nil
}
