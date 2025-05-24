package model

import (
	"time"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
)

type UnsavedUserAccessToken struct {
	User             User
	AccessToken      pkg.JwtToken
	SourceUpdateDate time.Time
	RegisteredDate   time.Time
	ExpireDate       *time.Time
}

type UserAccessToken struct {
	UserAccessTokenId uint
	UnsavedUserAccessToken
}

func CreateUserAccessToken(user User, accessToken pkg.JwtToken, registerDate time.Time, expireDate time.Time) UnsavedUserAccessToken {
	return UnsavedUserAccessToken{
		User:            user,
		AccessToken:     accessToken,
		SourceUpdateDate user.UpdateDate,
		RegisterDate:    registerDate,
		ExpireDate:      expireDate,
	}
}

func NewUserAccessToken(userAccessTokenId uint, user User, accessToken pkg.JwtToken, sourceUpdateDate time.Time, registerDate time.Time, expireDate *time.Time) UserRefreshToken {
	return UserRefreshToken{
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
