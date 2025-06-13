package user

import (
	text "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type UnsavedUserAccessToken struct {
	User             User
	AccessToken      text.JwtToken
	SourceUpdateDate time.Time
	RegisterDate     time.Time
	ExpireDate       time.Time
}

type UserAccessToken struct {
	PersistKey uint
	UnsavedUserAccessToken
}

func CreateUserAccessToken(
	user User,
	accessToken text.JwtToken,
	registerDate time.Time,
	expireDate time.Time,
) UnsavedUserAccessToken {
	return UnsavedUserAccessToken{
		User:             user,
		AccessToken:      accessToken,
		SourceUpdateDate: user.UpdateDate,
		RegisterDate:     registerDate,
		ExpireDate:       expireDate,
	}
}

func NewUserAccessToken(
	persistKey uint,
	user User,
	accessToken text.JwtToken,
	sourceUpdateDate time.Time,
	registerDate time.Time,
	expireDate time.Time,
) UserAccessToken {
	return UserAccessToken{
		PersistKey: persistKey,
		UnsavedUserAccessToken: UnsavedUserAccessToken{
			User:             user,
			AccessToken:      accessToken,
			SourceUpdateDate: sourceUpdateDate,
			RegisterDate:     registerDate,
			ExpireDate:       expireDate,
		},
	}
}
