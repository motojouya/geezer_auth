package user

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
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
	common.RequestHeader
	UserVerifyEmail
}

func (u UserVerifyEmailRequest) GetEmail() (pkgText.Email, error) {
	return pkgText.NewEmail(u.Email)
}

func (u UserVerifyEmailRequest) GetVerifyToken() (text.Token, error) {
	return text.NewToken(u.VerifyToken)
}
