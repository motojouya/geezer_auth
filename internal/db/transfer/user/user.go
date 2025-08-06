package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	"github.com/motojouya/geezer_auth/pkg/shelter/text"
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
	dbMap.AddTableWithName(User{}, "users").SetKeys(true, "PersistKey")
}

var SelectUser = utility.Dialect.From(goqu.T("users").As("u")).Select(
	goqu.I("u.persist_key").As("persist_key"),
	goqu.I("u.identifier").As("identifier"),
	goqu.I("u.email_identifier").As("email_identifier"),
	goqu.I("u.name").As("name"),
	goqu.I("u.bot_flag").As("bot_flag"),
	goqu.I("u.register_date").As("register_date"),
	goqu.I("u.update_date").As("update_date"),
)

func FromShelterUser(shelterUser shelter.UnsavedUser) User {
	return User{
		Identifier:     string(shelterUser.Identifier),
		ExposeEmailId:  string(shelterUser.ExposeEmailId),
		Name:           string(shelterUser.Name),
		BotFlag:        shelterUser.BotFlag,
		RegisteredDate: shelterUser.RegisteredDate,
		UpdateDate:     shelterUser.UpdateDate,
	}
}

func (u User) ToShelterUser() (shelter.User, error) {
	var identifier, idErr = text.NewIdentifier(u.Identifier)
	if idErr != nil {
		return shelter.User{}, idErr
	}

	var emailId, emailErr = text.NewEmail(u.ExposeEmailId)
	if emailErr != nil {
		return shelter.User{}, emailErr
	}

	var name, nameErr = text.NewName(u.Name)
	if nameErr != nil {
		return shelter.User{}, nameErr
	}

	return shelter.NewUser(
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
