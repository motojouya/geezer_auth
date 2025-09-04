package user

import (
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type EmailGetter interface {
	GetEmail() (pkgText.Email, error)
}

type UserChangeEmail struct {
	Email string `json:"email"`
}

type UserChangeEmailRequest struct {
	UserChangeEmail
}

func (u UserChangeEmailRequest) GetEmail() (pkgText.Email, error) {
	return pkgText.NewEmail(u.Email)
}
