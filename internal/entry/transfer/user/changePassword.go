package user

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	"github.com/motojouya/geezer_auth/internal/shelter/text"
)

type PasswordGetter interface {
	GetPassword() (text.Password, error)
}

type UserChangePassword struct {
	Password string `json:"password"`
}

type UserChangePasswordRequest struct {
	common.RequestHeader
	UserChangePassword
}

func (u UserChangePasswordRequest) GetPassword() (text.Password, error) {
	return text.NewPassword(u.Password)
}
