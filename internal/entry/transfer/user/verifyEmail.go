package user

import (
	text "github.com/motojouya/geezer_auth/internal/shelter/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type EmailVerifier interface {
	EmailGetter
	GetVerifyToken() (text.Token, error)
}

type UserVerifyEmail struct {
	Email       string `json:"email"`
	VerifyToken string `json:"verify_token"`
}

type UserVerifyEmailRequest struct {
	UserVerifyEmail UserVerifyEmail `http:"body"`
}

func (u UserVerifyEmailRequest) GetEmail() (pkgText.Email, error) {
	return pkgText.NewEmail(u.UserVerifyEmail.Email)
}

func (u UserVerifyEmailRequest) GetVerifyToken() (text.Token, error) {
	return text.NewToken(u.UserVerifyEmail.VerifyToken)
}
