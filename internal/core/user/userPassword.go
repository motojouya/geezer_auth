package user

import (
	"time"
	"github.com/motojouya/geezer_auth/internal/core/text"
)

type UnsavedUserPassword struct {
	User           User
	Password       text.HashedPassword
	RegisteredDate time.Time
	ExpireDate     *time.Time
}

type UserPassword struct {
	UserPasswordID uint
	UnsavedUserPassword
}

func CreateUserPassword(
	user User,
	password text.HashedPassword,
	registeredDate time.Time
) *UnsavedUserPassword {
	return &UnsavedUserPassword{
		User:           user,
		Password:       password,
		RegisteredDate: registeredDate,
		ExpireDate:     nil,
	}
}

func NewUserPassword(
	userPasswordID uint,
	user User,
	password text.HashedPassword,
	registeredDate time.Time,
	expireDate time.Time
) *UserPassword {
	return &UserPassword{
		UserPasswordID userPasswordID
		UnsavedUserPassword: UnsavedUserPassword{
			User:           user,
			Password:       password,
			RegisteredDate: registeredDate,
			ExpireDate:     expireDate,
		}
	}
}
