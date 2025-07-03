package user

import (
	core "github.com/motojouya/geezer_auth/internal/core/user"
	text "github.com/motojouya/geezer_auth/internal/core/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type UserVerifyEmail struct {
    Email       string `json:"email"`
    VerifyToken string `json:"verify_token"`
}

type UserVerifyEmailRequest struct {
	UserVerifyEmail UserVerifyEmail `http:"body"`
}

func (u UserVerifyEmailRequest) ToCoreUserEmail(user core.User, registerDate time.Time) (*core.UnsavedUserEmail, error) {
	var email, emailErr = pkgText.NewEmail(u.UserVerifyEmail.Email)
	if emailErr != nil {
		return &core.UnsavedUserEmail{}, emailErr
	}

	var verifyToken, tokenErr = text.NewToken(u.UserVerifyEmail.VerifyToken)
	if tokenErr != nil {
		return &core.UnsavedUserEmail{}, tokenErr
	}

	return core.CreateUserEmail(user, email, verifyToken, registerDate), nil
}
