package user

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	text "github.com/motojouya/geezer_auth/internal/core/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type UserRegister struct {
    Email     string `json:"email"`
    Name      string `json:"name"`
    Bot       bool   `json:"bot"`
    Password  string `json:"password"`
}

type UserRegisterRequest struct {
	UserRegister UserRegister `http:"body"`
}

type UserRegisterResponse struct {
    User common.User `json:"user"`
    RefreshToken string `json:"refresh_token"`
    AccessToken string `json:"access_token"`
}

func FromCoreUserAuthentic(coreUser *core.UserAuthentic, refreshToken text.Token, accessToken pkgText.JwtToken) *UserRegisterResponse {
	var commonUser = common.FromCoreUser(coreUser)
	return &UserRegisterResponse{
		User: *commonUser,
		RefreshToken: string(refreshToken),
		AccessToken: string(accessToken),
	}
}

func (u UserRegisterRequest) ToCoreUser(identifier pkgText.Identifier, registerDate time.Time) (core.UnsavedUser, error) {
	var emailId, emailErr = pkgText.NewEmail(u.UserRegister.Email)
	if emailErr != nil {
		return core.UnsavedUser{}, emailErr
	}

	var name, nameErr = pkgText.NewName(u.UserRegister.Name)
	if nameErr != nil {
		return core.UnsavedUser{}, nameErr
	}

	return core.CreateUser(identifier, emailId, name, u.UserRegister.Bot, registerDate), nil
}

func (u UserRegisterRequest) GetPassword() (text.Password, error) {
	return text.NewPassword(u.UserRegister.Password)
}

// TODO 以下は別ファイル。requestとresponseはそもそも別ファイルのがいいかも
type UserUpdateResponse struct {
    User common.User `json:"user"`
    AccessToken string `json:"access_token"`
}

type UserGetResponse struct {
    User common.User `json:"user"`
}
