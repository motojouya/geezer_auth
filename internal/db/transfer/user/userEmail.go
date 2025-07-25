package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	text "github.com/motojouya/geezer_auth/internal/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	pkg "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type UserEmail struct {
	PersistKey     uint       `db:"persist_key,primarykey,autoincrement"`
	UserPersistKey uint       `db:"user_persist_key"`
	Email          string     `db:"email"`
	VerifyToken    string     `db:"verify_token"`
	RegisterDate   time.Time  `db:"register_date"`
	VerifyDate     *time.Time `db:"verify_date"`
	ExpireDate     *time.Time `db:"expire_date"`
}

type UserEmailFull struct {
	UserEmail
	UserIdentifier     string    `db:"user_identifier"`
	UserExposeEmailId  string    `db:"user_email_identifier"`
	UserName           string    `db:"user_name"`
	UserBotFlag        bool      `db:"user_bot_flag"`
	UserRegisteredDate time.Time `db:"user_register_date"`
	UserUpdateDate     time.Time `db:"user_update_date"`
}

func AddUserEmailTable(dbMap *gorp.DbMap) {
	dbMap.AddTableWithName(UserEmail{}, "user_email").SetKeys(true, "PersistKey")
}

// var SelectUserEmail = utility.Dialect.From("user_email").As("ue").Select(
// 	goqu.C("ue.persist_key").As("persist_key"),
// 	goqu.C("ue.user_persist_key").As("user_persist_key"),
// 	goqu.C("ue.email").As("email"),
// 	goqu.C("ue.verify_token").As("verify_token"),
// 	goqu.C("ue.register_date").As("register_date"),
// 	goqu.C("ue.verify_date").As("verify_date"),
// 	goqu.C("ue.expire_date").As("expire_date"),
// )

var SelectUserEmail = utility.Dialect.From(goqu.T("user_email").As("ue")).InnerJoin(
	goqu.T("users").As("u"),
	goqu.On(goqu.Ex{"ue.user_persist_key": goqu.I("u.persist_key")}),
).Select(
	goqu.I("ue.persist_key").As("persist_key"),
	goqu.I("ue.user_persist_key").As("user_persist_key"),
	goqu.I("u.identifier").As("user_identifier"),
	goqu.I("u.email_identifier").As("user_email_identifier"),
	goqu.I("u.name").As("user_name"),
	goqu.I("u.bot_flag").As("user_bot_flag"),
	goqu.I("u.register_date").As("user_register_date"),
	goqu.I("u.update_date").As("user_update_date"),
	goqu.I("ue.email").As("email"),
	goqu.I("ue.verify_token").As("verify_token"),
	goqu.I("ue.register_date").As("register_date"),
	goqu.I("ue.verify_date").As("verify_date"),
	goqu.I("ue.expire_date").As("expire_date"),
)

func FromCoreUserEmail(coreUserEmail *core.UnsavedUserEmail) *UserEmail {
	return &UserEmail{
		UserPersistKey: coreUserEmail.User.PersistKey,
		Email:          string(coreUserEmail.Email),
		VerifyToken:    string(coreUserEmail.VerifyToken),
		RegisterDate:   coreUserEmail.RegisterDate,
		VerifyDate:     coreUserEmail.VerifyDate,
		ExpireDate:     coreUserEmail.ExpireDate,
	}
}

func (u UserEmailFull) ToCoreUserEmail() (*core.UserEmail, error) {
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
		return &core.UserEmail{}, userErr
	}

	var email, emailErr = pkg.NewEmail(u.Email)
	if emailErr != nil {
		return &core.UserEmail{}, emailErr
	}

	var verifyToken, tokenErr = text.NewToken(u.VerifyToken)
	if tokenErr != nil {
		return &core.UserEmail{}, tokenErr
	}

	return core.NewUserEmail(
		u.PersistKey,
		user,
		email,
		verifyToken,
		u.RegisterDate,
		u.VerifyDate,
		u.ExpireDate,
	), nil
}

// testdata投入時に楽するためのもの。アプリケーションからは利用を想定しない。
func NewUserEmail(persistKey uint, userPersistKey uint, email string, verifyToken string, registerDate time.Time, verifyDate *time.Time, expireDate *time.Time) *UserEmail {
	return &UserEmail{
		PersistKey:     persistKey,
		UserPersistKey: userPersistKey,
		Email:          email,
		VerifyToken:    verifyToken,
		RegisterDate:   registerDate,
		VerifyDate:     verifyDate,
		ExpireDate:     expireDate,
	}
}
