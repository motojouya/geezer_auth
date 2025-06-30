package user

import (
	"github.com/motojouya/geezer_auth/internal/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"time"
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
