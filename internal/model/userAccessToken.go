package model

import (
	"time"
)

type UnsavedUserAccessToken struct {
	User             User
	AccessToken      AccessToken
	SourceUpdateDate time.Time
	RegisteredDate   time.Time
	ExpireDate       *time.Time
}

type UserAccessToken struct {
	UserAccessTokenId uint
	UnsavedUserAccessToken
}

func CreateUserAccessToken(user User, accessToken AccessToken, registerDate time.Time) UnsavedUserAccessToken {
	return UnsavedUserAccessToken{
		User:            user,
		AccessToken:     accessToken,
		SourceUpdateDate user.UpdateDate,
		RegisterDate:    registerDate,
		ExpireDate:      nil,
	}
}

func NewUserAccessToken(userAccessTokenId uint, user User, accessToken AccessToken, sourceUpdateDate time.Time, registerDate time.Time, expireDate *time.Time) UserRefreshToken {
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
