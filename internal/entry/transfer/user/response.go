package user

import (
	text "github.com/motojouya/geezer_auth/internal/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
)

type UserGetResponse struct {
	User common.User `json:"user"`
}

type UserUpdateResponse struct {
	UserGetResponse
	AccessToken string `json:"access_token"`
}

type UserRegisterResponse struct {
	UserUpdateResponse
	RefreshToken string `json:"refresh_token"`
}

func FromCoreUserAuthenticToGetResponse(coreUser *core.UserAuthentic) *UserGetResponse {
	var commonUser = common.FromCoreUser(coreUser)
	return &UserGetResponse{
		User: *commonUser,
	}
}

func FromCoreUserAuthenticToUpdateResponse(coreUser *core.UserAuthentic, accessToken pkgText.JwtToken) *UserUpdateResponse {
	var userGetResponse = FromCoreUserAuthenticToGetResponse(coreUser)
	return &UserUpdateResponse{
		UserGetResponse: *userGetResponse,
		AccessToken:     string(accessToken),
	}
}

func FromCoreUserAuthenticToRegisterResponse(coreUser *core.UserAuthentic, refreshToken text.Token, accessToken pkgText.JwtToken) *UserRegisterResponse {
	var userUpdateResponse = FromCoreUserAuthenticToUpdateResponse(coreUser, accessToken)
	return &UserRegisterResponse{
		UserUpdateResponse: *userUpdateResponse,
		RefreshToken:       string(refreshToken),
	}
}
