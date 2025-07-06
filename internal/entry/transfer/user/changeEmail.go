package user

import (
	text "github.com/motojouya/geezer_auth/internal/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type UserChangeEmail struct {
	Email string `json:"email"`
}

type UserChangeEmailRequest struct {
	UserChangeEmail UserChangeEmail `http:"body"`
}

func (u UserChangeEmailRequest) ToCoreUserEmail(user core.User, verifyToken text.Token, registerDate time.Time) (*core.UnsavedUserEmail, error) {
	var email, emailErr = pkgText.NewEmail(u.UserChangeEmail.Email)
	if emailErr != nil {
		return &core.UnsavedUserEmail{}, emailErr
	}

	return core.CreateUserEmail(user, email, verifyToken, registerDate), nil
}
