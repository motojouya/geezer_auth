package model

import (
	"time"
	"github.com/google/uuid"
)

type UnsavedUserRefreshToken struct {
	User               User
	RefreshToken       uuid.UUID
	RegisteredDate     time.Time
	ExpireDate         time.Time
}

type UserRefreshToken struct {
	UserRefreshTokenId uint
	UnsavedUserRefreshToken
}

// FIXME 外から環境変数で設定できてもいいかも
const TokenValidityPeriodDays = 50

func CreateUserRefreshToken(user User, refreshToken uuid.UUID, registerDate time.Time) UnsavedUserRefreshToken {
	var expireDate = registerDate.Add(TokenValidityPeriodDays * time.Day)

	return UnsavedUserRefreshToken{
		User:         user,
		RefreshToken: refreshToken,
		RegisterDate: registerDate,
		ExpireDate:   expireDate,
	}
}

func NewUserRefreshToken(userRefreshTokenId uint, user User, refreshToken uuid.UUID,registerDate time.Time, expireDate time.Time) UserRefreshToken {
	return UserRefreshToken{
		UserRefreshTokenId: userRefreshTokenId,
		UnsavedUserRefreshToken: UnsavedUserRefreshToken{
			User:         user,
			RefreshToken: refreshToken,
			RegisterDate: registerDate,
			ExpireDate:   expireDate,
		},
	}
}
