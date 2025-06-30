package user

import (
	"github.com/motojouya/geezer_auth/internal/core/text"
	"time"
)

type UnsavedUserPassword struct {
	User           User
	Password       text.HashedPassword
	RegisteredDate time.Time
	ExpireDate     *time.Time
}

type UserPassword struct {
	PersistKey uint
	UnsavedUserPassword
}

func CreateUserPassword(
	user User,
	password text.HashedPassword,
	registeredDate time.Time,
) *UnsavedUserPassword {
	return &UnsavedUserPassword{
		User:           user,
		Password:       password,
		RegisteredDate: registeredDate,
		ExpireDate:     nil,
	}
}

func NewUserPassword(
	persistKey uint,
	user User,
	password text.HashedPassword,
	registeredDate time.Time,
	expireDate *time.Time,
) *UserPassword {
	return &UserPassword{
		PersistKey: persistKey,
		UnsavedUserPassword: UnsavedUserPassword{
			User:           user,
			Password:       password,
			RegisteredDate: registeredDate,
			ExpireDate:     expireDate,
		},
	}
}
