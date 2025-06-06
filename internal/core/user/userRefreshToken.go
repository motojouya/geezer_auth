package user

import (
	"time"
	"github.com/motojouya/geezer_auth/internal/core/text"
)

type UnsavedUserRefreshToken struct {
	User               User
	RefreshToken       text.Token
	RegisteredDate     time.Time
	ExpireDate         time.Time
}

type UserRefreshToken struct {
	PersistKey uint
	UnsavedUserRefreshToken
}

// FIXME 外から環境変数で設定できてもいいかも
const TokenValidityPeriodDays = 50

func CreateUserRefreshToken(
	user User,
	refreshToken text.Token,
	registerDate time.Time
) UnsavedUserRefreshToken {
	var expireDate = registerDate.Add(TokenValidityPeriodDays * time.Day)

	return UnsavedUserRefreshToken{
		User:         user,
		RefreshToken: refreshToken,
		RegisterDate: registerDate,
		ExpireDate:   expireDate,
	}
}

func NewUserRefreshToken(
	persistKey uint
	user User,
	refreshToken text.Token,
	registerDate time.Time,
	expireDate time.Time
) UserRefreshToken {
	return UserRefreshToken{
		PersistKey:              persistKey,
		UnsavedUserRefreshToken: UnsavedUserRefreshToken{
			User:         user,
			RefreshToken: refreshToken,
			RegisterDate: registerDate,
			ExpireDate:   expireDate,
		},
	}
}
