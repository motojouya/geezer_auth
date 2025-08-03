package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/shelter/text"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"time"
)

type UserRefreshToken struct {
	PersistKey     uint      `db:"persist_key,primarykey,autoincrement"`
	UserPersistKey uint      `db:"user_persist_key"`
	RefreshToken   string    `db:"refresh_token"`
	RegisterDate   time.Time `db:"register_date"`
	ExpireDate     time.Time `db:"expire_date"`
}

type UserRefreshTokenFull struct {
	UserRefreshToken
	UserIdentifier     string    `db:"user_identifier"`
	UserExposeEmailId  string    `db:"user_email_identifier"`
	UserName           string    `db:"user_name"`
	UserBotFlag        bool      `db:"user_bot_flag"`
	UserRegisteredDate time.Time `db:"user_register_date"`
	UserUpdateDate     time.Time `db:"user_update_date"`
}

func AddUserRefreshTokenTable(dbMap *gorp.DbMap) {
	dbMap.AddTableWithName(UserRefreshToken{}, "user_refresh_token").SetKeys(true, "PersistKey")
}

// var SelectUserRefreshToken = utility.Dialect.From("user_refresh_token").As("urt").Select(
// 	goqu.C("urt.persist_key").As("persist_key"),
// 	goqu.C("urt.user_persist_key").As("user_persist_key"),
// 	goqu.C("urt.refresh_token").As("refresh_token"),
// 	goqu.C("urt.register_date").As("register_date"),
// 	goqu.C("urt.expire_date").As("expire_date"),
// )

var SelectUserRefreshToken = utility.Dialect.From(goqu.T("user_refresh_token").As("urt")).InnerJoin(
	goqu.T("users").As("u"),
	goqu.On(goqu.Ex{"urt.user_persist_key": goqu.I("u.persist_key")}),
).Select(
	goqu.I("urt.persist_key").As("persist_key"),
	goqu.I("urt.user_persist_key").As("user_persist_key"),
	goqu.I("u.identifier").As("user_identifier"),
	goqu.I("u.email_identifier").As("user_email_identifier"),
	goqu.I("u.name").As("user_name"),
	goqu.I("u.bot_flag").As("user_bot_flag"),
	goqu.I("u.register_date").As("user_register_date"),
	goqu.I("u.update_date").As("user_update_date"),
	goqu.I("urt.refresh_token").As("refresh_token"),
	goqu.I("urt.register_date").As("register_date"),
	goqu.I("urt.expire_date").As("expire_date"),
)

func FromCoreUserRefreshToken(shelterUserRefreshToken shelter.UnsavedUserRefreshToken) UserRefreshToken {
	return UserRefreshToken{
		UserPersistKey: shelterUserRefreshToken.User.PersistKey,
		RefreshToken:   string(shelterUserRefreshToken.RefreshToken),
		RegisterDate:   shelterUserRefreshToken.RegisterDate,
		ExpireDate:     shelterUserRefreshToken.ExpireDate,
	}
}

func (u UserRefreshTokenFull) ToCoreUserRefreshToken() (shelter.UserRefreshToken, error) {
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
		return shelter.UserRefreshToken{}, userErr
	}

	var token, tokenErr = text.NewToken(u.RefreshToken)
	if tokenErr != nil {
		return shelter.UserRefreshToken{}, tokenErr
	}

	return shelter.NewUserRefreshToken(
		u.PersistKey,
		user,
		token,
		u.RegisterDate,
		u.ExpireDate,
	), nil
}

// testdata投入時に楽するためのもの。アプリケーションからは利用を想定しない。
func NewUserRefreshToken(persistKey uint, userPersistKey uint, refreshToken string, registerDate time.Time, expireDate time.Time) UserRefreshToken {
	return UserRefreshToken{
		PersistKey:     persistKey,
		UserPersistKey: userPersistKey,
		RefreshToken:   refreshToken,
		RegisterDate:   registerDate,
		ExpireDate:     expireDate,
	}
}
