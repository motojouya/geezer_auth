package user

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	text "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
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

func FromCoreUserAuthenticToGetResponse(shelterUser *shelter.UserAuthentic) *UserGetResponse {
	var commonUser = common.FromCoreUser(shelterUser)
	return &UserGetResponse{
		User: *commonUser,
	}
}

func FromCoreUserAuthenticToUpdateResponse(shelterUser *shelter.UserAuthentic, accessToken pkgText.JwtToken) *UserUpdateResponse {
	var userGetResponse = FromCoreUserAuthenticToGetResponse(shelterUser)
	return &UserUpdateResponse{
		UserGetResponse: *userGetResponse,
		AccessToken:     string(accessToken),
	}
}

func FromCoreUserAuthenticToRegisterResponse(shelterUser *shelter.UserAuthentic, refreshToken text.Token, accessToken pkgText.JwtToken) *UserRegisterResponse {
	var userUpdateResponse = FromCoreUserAuthenticToUpdateResponse(shelterUser, accessToken)
	return &UserRegisterResponse{
		UserUpdateResponse: *userUpdateResponse,
		RefreshToken:       string(refreshToken),
	}
}
