package user

import (
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
)

type EmailGetter interface {
	GetEmail() (pkgText.Email, error)
}

type UserChangeEmail struct {
	Email string `json:"email"`
}

type UserChangeEmailRequest struct {
	common.RequestHeader
	UserChangeEmail
}

func (u UserChangeEmailRequest) GetEmail() (pkgText.Email, error) {
	return pkgText.NewEmail(u.Email)
}
