package model

import (
	"time"
)

type UnsavedUserEmail struct {
	User         User
	Email        Email
	Verified     bool
	RegisterDate time.Time
	ExpireDate   *time.Time
}

type UserEmail struct {
	UserEmailID    uint
	UnsavedUserEmail
}

func CreateUserEmail(user User, email Email, registerDate time.Time) UnsavedUserEmail {
	return UnsavedUserEmail{
		User:         user,
		Email:        email,
		Verified:     false,
		RegisterDate: registerDate,
		ExpireDate:   nil,
	}
}

func NewUserEmail(userEmailId uint, user User, email Email, verified bool, registerDate time.Time, expireDate *time.Time) UserEmail {
	return UserEmail{
		UserEmailID: userEmailId,
		UnsavedUserEmail: UnsavedUserEmail{
			User:         user,
			Email:        email,
			Verified:     verified,
			RegisterDate: registerDate,
			ExpireDate:   expireDate,
		},
	}
}
