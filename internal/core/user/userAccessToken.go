package user

import (
	"time"
	text "github.com/motojouya/geezer_auth/pkg/core/text"
)

type UnsavedUserAccessToken struct {
	User             User
	AccessToken      text.JwtToken
	SourceUpdateDate time.Time
	RegisteredDate   time.Time
	ExpireDate       *time.Time
}

type UserAccessToken struct {
	UserAccessTokenId uint
	UnsavedUserAccessToken
}

func CreateUserAccessToken(
	user User,
	accessToken text.JwtToken,
	registerDate time.Time,
	expireDate time.Time
) *UnsavedUserAccessToken {
	return &UnsavedUserAccessToken{
		User:            user,
		AccessToken:     accessToken,
		SourceUpdateDate user.UpdateDate,
		RegisterDate:    registerDate,
		ExpireDate:      expireDate,
	}
}

func NewUserAccessToken(
	userAccessTokenId uint,
	user User,
	accessToken text.JwtToken,
	sourceUpdateDate time.Time,
	registerDate time.Time,
	expireDate *time.Time
) *UserRefreshToken {
	return &UserRefreshToken{
		UserAccessTokenId: userAccessTokenId,
		UnsavedUserAccessToken: UnsavedUserAccessToken{
			User:            user,
			AccessToken:     accessToken,
			SourceUpdateDate sourceUpdateDate,
			RegisterDate:    registerDate,
			ExpireDate:      expireDate,
		}
	}
}
