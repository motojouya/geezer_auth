package user

import (
	text "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"time"
)

type UserChangeEmail struct {
	Email string `json:"email"`
}

type UserChangeEmailRequest struct {
	UserChangeEmail UserChangeEmail `http:"body"`
}

func (u UserChangeEmailRequest) ToCoreUserEmail(user shelter.User, verifyToken text.Token, registerDate time.Time) (*shelter.UnsavedUserEmail, error) {
	var email, emailErr = pkgText.NewEmail(u.UserChangeEmail.Email)
	if emailErr != nil {
		return &shelter.UnsavedUserEmail{}, emailErr
	}

	return shelter.CreateUserEmail(user, email, verifyToken, registerDate), nil
}
