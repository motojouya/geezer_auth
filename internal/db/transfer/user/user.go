package user

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"time"
)

type User struct {
	PersistKey     uint
	Identifier     string
	ExposeEmailId  string
	Name           string
	BotFlag        bool
	RegisteredDate time.Time
	UpdateDate     time.Time
}

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
	var identifier, idErr = text.NewIdentifier(u.Identifier);
	if idErr != nil {
		return core.User{}, idErr
	}

	var emailId, emailErr = text.NewEmail(u.ExposeEmailId);
	if emailErr != nil {
		return core.User{}, emailErr
	}

	var name, nameErr = text.NewName(u.Name);
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
