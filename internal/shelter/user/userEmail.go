package user

import (
	text "github.com/motojouya/geezer_auth/internal/shelter/text"
	pkg "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"time"
)

/*
 * Userの関連情報でEmailの履歴
 * 履歴テーブルはmodel不要に感じるが、verifiedかの判定がありそうで、、modelにしている
 */
type UnsavedUserEmail struct {
	User         User
	Email        pkg.Email
	VerifyToken  text.Token
	RegisterDate time.Time
	VerifyDate   *time.Time
	ExpireDate   *time.Time
}

type UserEmail struct {
	PersistKey uint
	UnsavedUserEmail
}

func CreateUserEmail(
	user User,
	email pkg.Email,
	verifyToken text.Token,
	registerDate time.Time,
) *UnsavedUserEmail {
	return &UnsavedUserEmail{
		User:         user,
		Email:        email,
		VerifyToken:  verifyToken,
		RegisterDate: registerDate,
		VerifyDate:   nil,
		ExpireDate:   nil,
	}
}

func NewUserEmail(
	persistKey uint,
	user User,
	email pkg.Email,
	verifyToken text.Token,
	registerDate time.Time,
	verifyDate *time.Time,
	expireDate *time.Time,
) *UserEmail {
	return &UserEmail{
		PersistKey: persistKey,
		UnsavedUserEmail: UnsavedUserEmail{
			User:         user,
			Email:        email,
			VerifyToken:  verifyToken,
			RegisterDate: registerDate,
			VerifyDate:   verifyDate,
			ExpireDate:   expireDate,
		},
	}
}
