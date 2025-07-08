package user

import (
	"github.com/motojouya/geezer_auth/internal/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"time"
	"github.com/motojouya/geezer_auth/internal/db"
	"github.com/doug-martin/goqu/v9"
)

type UserPassword struct {
	PersistKey     uint       `db:"persist_key"`
	UserPersistKey uint       `db:"user_persist_key"`
	Password       string     `db:"password"`
	RegisteredDate time.Time  `db:"register_date"`
	ExpireDate     *time.Time `db:"expire_date"`
}

type UserPasswordFull struct {
	UserPassword
	UserIdentifier     string    `db:"user_identifier"`
	UserExposeEmailId  string    `db:"user_email_identifier"`
	UserName           string    `db:"user_name"`
	UserBotFlag        bool      `db:"user_bot_flag"`
	UserRegisteredDate time.Time `db:"user_register_date"`
	UserUpdateDate     time.Time `db:"user_update_date"`
}

var SelectUserPassword = db.Dialect.From("user_password").As("up").Select(
	goqu.C("up.persist_key").As("persist_key"),
	goqu.C("up.user_persist_key").As("user_persist_key"),
	goqu.C("up.password").As("password"),
	goqu.C("up.register_date").As("register_date"),
	goqu.C("up.expire_date").As("expire_date"),
)

var SelectFullUserPassword = db.Dialect.From("user_password").As("up").InnerJoin(
	goqu.T("user").As("u"),
	goqu.On(goqu.Ex{"uat.user_persist_key": goqu.I("u.persist_key")}),
).Select(
	goqu.C("up.persist_key").As("persist_key"),
	goqu.C("up.user_persist_key").As("user_persist_key"),
	goqu.C("u.identifier").As("user_identifier"),
	goqu.C("u.email_identifier").As("user_email_identifier"),
	goqu.C("u.name").As("user_name"),
	goqu.C("u.bot_flag").As("user_bot_flag"),
	goqu.C("u.register_date").As("user_register_date"),
	goqu.C("u.update_date").As("user_update_date"),
	goqu.C("up.password").As("password"),
	goqu.C("up.register_date").As("register_date"),
	goqu.C("up.expire_date").As("expire_date"),
)

func FromCoreUserPassword(coreUserPassword *core.UnsavedUserPassword) *UserPassword {
	return &UserPassword{
		UserPersistKey: coreUserPassword.User.PersistKey,
		Password:       string(coreUserPassword.Password),
		RegisteredDate: coreUserPassword.RegisteredDate,
		ExpireDate:     coreUserPassword.ExpireDate,
	}
}

func (u UserPasswordFull) ToCoreUserPassword() (*core.UserPassword, error) {
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
		return &core.UserPassword{}, userErr
	}

	return core.NewUserPassword(
		u.PersistKey,
		user,
		text.NewHashedPassword(u.Password),
		u.RegisteredDate,
		u.ExpireDate,
	), nil
}
