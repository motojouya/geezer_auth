package user

import (
	"github.com/motojouya/geezer_auth/internal/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"time"
)

type UserPassword struct {
	PersistKey     uint
	UserPersistKey uint
	Password       string
	RegisteredDate time.Time
	ExpireDate     *time.Time
}

type UserPasswordFull struct {
	UserPassword
	UserIdentifier     string
	UserExposeEmailId  string
	UserName           string
	UserBotFlag        bool
	UserRegisteredDate time.Time
	UserUpdateDate     time.Time
}

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
