package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"github.com/motojouya/geezer_auth/internal/shelter/text"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	"time"
)

type UserPassword struct {
	PersistKey     uint       `db:"persist_key,primarykey,autoincrement"`
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

func AddUserPasswordTable(dbMap *gorp.DbMap) {
	dbMap.AddTableWithName(UserPassword{}, "user_password").SetKeys(true, "PersistKey")
}

// var SelectUserPassword = utility.Dialect.From("user_password").As("up").Select(
// 	goqu.C("up.persist_key").As("persist_key"),
// 	goqu.C("up.user_persist_key").As("user_persist_key"),
// 	goqu.C("up.password").As("password"),
// 	goqu.C("up.register_date").As("register_date"),
// 	goqu.C("up.expire_date").As("expire_date"),
// )

var SelectUserPassword = utility.Dialect.From(goqu.T("user_password").As("up")).InnerJoin(
	goqu.T("users").As("u"),
	goqu.On(goqu.Ex{"up.user_persist_key": goqu.I("u.persist_key")}),
).Select(
	goqu.I("up.persist_key").As("persist_key"),
	goqu.I("up.user_persist_key").As("user_persist_key"),
	goqu.I("u.identifier").As("user_identifier"),
	goqu.I("u.email_identifier").As("user_email_identifier"),
	goqu.I("u.name").As("user_name"),
	goqu.I("u.bot_flag").As("user_bot_flag"),
	goqu.I("u.register_date").As("user_register_date"),
	goqu.I("u.update_date").As("user_update_date"),
	goqu.I("up.password").As("password"),
	goqu.I("up.register_date").As("register_date"),
	goqu.I("up.expire_date").As("expire_date"),
)

func FromCoreUserPassword(shelterUserPassword *shelter.UnsavedUserPassword) *UserPassword {
	return &UserPassword{
		UserPersistKey: shelterUserPassword.User.PersistKey,
		Password:       string(shelterUserPassword.Password),
		RegisteredDate: shelterUserPassword.RegisteredDate,
		ExpireDate:     shelterUserPassword.ExpireDate,
	}
}

func (u UserPasswordFull) ToCoreUserPassword() (*shelter.UserPassword, error) {
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
		return &shelter.UserPassword{}, userErr
	}

	return shelter.NewUserPassword(
		u.PersistKey,
		user,
		text.NewHashedPassword(u.Password),
		u.RegisteredDate,
		u.ExpireDate,
	), nil
}

// testdata投入時に楽するためのもの。アプリケーションからは利用を想定しない。
func NewUserPassword(persistKey uint, userPersistKey uint, password string, registeredDate time.Time, expireDate *time.Time) *UserPassword {
	return &UserPassword{
		PersistKey:     persistKey,
		UserPersistKey: userPersistKey,
		Password:       password,
		RegisteredDate: registeredDate,
		ExpireDate:     expireDate,
	}
}
