package model

import (
	"time"
	"github.com/google/uuid"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
)

/*
 * Userの関連情報でEmailの履歴
 * 履歴テーブルはmodel不要に感じるが、verifiedかの判定がありそうで、、modelにしている
 */
type UnsavedUserEmail struct {
	User         User
	Email        pkg.Email
	VerifyToken  uuid.UUID
	RegisterDate time.Time
	VerifyDate   *time.Time
	ExpireDate   *time.Time
}

type UserEmail struct {
	UserEmailID    uint
	UnsavedUserEmail
}

func CreateUserEmail(user User, email pkg.Email, verifyToken uuid.UUID, registerDate time.Time) UnsavedUserEmail {
	return UnsavedUserEmail{
		User:         user,
		Email:        email,
		VerifyToken:  verifyToken,
		RegisterDate: registerDate,
		VerifyDate:   nil,
		ExpireDate:   nil,
	}
}

func NewUserEmail(userEmailId uint, user User, email pkg.Email, verifyToken uuid.UUID, registerDate time.Time, verifyDate time.Time, expireDate *time.Time) UserEmail {
	return UserEmail{
		UserEmailID: userEmailId,
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
