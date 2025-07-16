package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type User struct {
	PersistKey     uint      `db:"persist_key,primarykey,autoincrement"`
	Identifier     string    `db:"identifier"`
	ExposeEmailId  string    `db:"email_identifier"`
	Name           string    `db:"name"`
	BotFlag        bool      `db:"bot_flag"`
	RegisteredDate time.Time `db:"register_date"`
	UpdateDate     time.Time `db:"update_date"`
}

func AddUserTable(dbMap *gorp.DbMap) {
	dbMap.AddTableWithName(User{}, "user").SetKeys(true, "PersistKey")
}

var SelectUser = utility.Dialect.From("user").As("u").Select(
	goqu.C("u.persist_key").As("persist_key"),
	goqu.C("u.identifier").As("identifier"),
	goqu.C("u.email_identifier").As("email_identifier"),
	goqu.C("u.name").As("name"),
	goqu.C("u.bot_flag").As("bot_flag"),
	goqu.C("u.register_date").As("register_date"),
	goqu.C("u.update_date").As("update_date"),
)

func FromCoreUser(coreUser core.UnsavedUser) User {
	return User{
		Identifier:     string(coreUser.Identifier),
		ExposeEmailId:  string(coreUser.ExposeEmailId),
		Name:           string(coreUser.Name),
		BotFlag:        coreUser.BotFlag,
		RegisteredDate: coreUser.RegisteredDate,
		UpdateDate:     coreUser.UpdateDate,
	}
}

func (u User) ToCoreUser() (core.User, error) {
	var identifier, idErr = text.NewIdentifier(u.Identifier)
	if idErr != nil {
		return core.User{}, idErr
	}

	var emailId, emailErr = text.NewEmail(u.ExposeEmailId)
	if emailErr != nil {
		return core.User{}, emailErr
	}

	var name, nameErr = text.NewName(u.Name)
	if nameErr != nil {
		return core.User{}, nameErr
	}

	return core.NewUser(
		u.PersistKey,
		identifier,
		name,
		emailId,
		u.BotFlag,
		u.RegisteredDate,
		u.UpdateDate,
	), nil
}

// testdata投入時に楽するためのもの。アプリケーションからは利用を想定しない。
func NewUser(persistKey uint, identifier string, exposeEmailId string, name string, botFlag bool, registeredDate time.Time, updateDate time.Time) User {
	return User{
		PersistKey:     persistKey,
		Identifier:     identifier,
		ExposeEmailId:  exposeEmailId,
		Name:           name,
		BotFlag:        botFlag,
		RegisteredDate: registeredDate,
		UpdateDate:     updateDate,
	}
}
