package user

import (
	"github.com/motojouya/geezer_auth/internal/shelter/text"
	"time"
)

type UnsavedUserRefreshToken struct {
	User         User
	RefreshToken text.Token
	RegisterDate time.Time
	ExpireDate   time.Time
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
	registerDate time.Time,
) UnsavedUserRefreshToken {
	var expireDate = registerDate.AddDate(0, 0, TokenValidityPeriodDays)

	return UnsavedUserRefreshToken{
		User:         user,
		RefreshToken: refreshToken,
		RegisterDate: registerDate,
		ExpireDate:   expireDate,
	}
}

func NewUserRefreshToken(
	persistKey uint,
	user User,
	refreshToken text.Token,
	registerDate time.Time,
	expireDate time.Time,
) UserRefreshToken {
	return UserRefreshToken{
		PersistKey: persistKey,
		UnsavedUserRefreshToken: UnsavedUserRefreshToken{
			User:         user,
			RefreshToken: refreshToken,
			RegisterDate: registerDate,
			ExpireDate:   expireDate,
		},
	}
}
