package user

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"time"
)

type UserAccessToken struct {
	PersistKey       uint
	UserPersistKey   uint
	AccessToken      string
	SourceUpdateDate time.Time
	RegisterDate     time.Time
	ExpireDate       time.Time
}

type UserAccessTokenFull struct {
	PersistKey         uint
	UserPersistKey     uint
	UserIdentifier     string
	UserExposeEmailId  string
	UserName           string
	UserBotFlag        bool
	UserRegisteredDate time.Time
	UserUpdateDate     time.Time
	AccessToken        string
	SourceUpdateDate   time.Time
	RegisterDate       time.Time
	ExpireDate         time.Time
}

func FromCoreUserAccessToken(coreUserAccessToken core.UnsavedUserAccessToken) UserAccessToken {
	return UserAccessToken{
		UserPersistKey:   coreUserAccessToken.User.PersistKey,
		AccessToken:      string(coreUserAccessToken.AccessToken),
		SourceUpdateDate: coreUserAccessToken.SourceUpdateDate,
		RegisterDate:     coreUserAccessToken.RegisterDate,
		ExpireDate:       coreUserAccessToken.ExpireDate,
	}
}

func (ua UserAccessTokenFull) ToCoreUserAccessToken() (core.UserAccessToken, error) {
	var user, userErr = (User{
		PersistKey:     ua.UserPersistKey,
		Identifier:     ua.UserIdentifier,
		ExposeEmailId:  ua.UserExposeEmailId,
		Name:           ua.UserName,
		BotFlag:        ua.UserBotFlag,
		RegisteredDate: ua.UserRegisteredDate,
		UpdateDate:     ua.UserUpdateDate,
	}).ToCoreUser()
	if userErr != nil {
		return core.UserAccessToken{}, userErr
	}

	return core.NewUserAccessToken(
		ua.PersistKey,
		user,
		text.NewJwtToken(ua.AccessToken),
		ua.SourceUpdateDate,
		ua.RegisterDate,
		ua.ExpireDate,
	), nil
}
