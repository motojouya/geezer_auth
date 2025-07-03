package user

import (
	"github.com/motojouya/geezer_auth/internal/core/text"
)

type UserChangePassword struct {
    NextPassword string `json:"next_password"`
    NowPassword string `json:"now_password"`
}

type UserChangePasswordRequest struct {
	UserChangePassword UserChangePassword `http:"body"`
}

func (u UserChangePasswordRequest) GetNextPassword() (text.Password, error) {
	return text.NewPassword(u.UserChangePassword.NextPassword)
}

func (u UserChangePasswordRequest) GetNowPassword() (text.Password, error) {
	return text.NewPassword(u.UserChangePassword.NowPassword)
}
