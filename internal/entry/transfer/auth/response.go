package auth

import (
	text "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type AuthRefreshResponse struct {
	user.UserUpdateResponse
}

type AuthLoginResponse struct {
	user.UserRegisterResponse
}

func FromCoreUserAuthenticToRefreshResponse(shelterUser *shelter.UserAuthentic, accessToken pkgText.JwtToken) *AuthRefreshResponse {
	var userGetResponse = user.FromCoreUserAuthenticToGetResponse(shelterUser)
	return &AuthRefreshResponse{
		UserUpdateResponse: user.UserUpdateResponse{
			UserGetResponse: *userGetResponse,
			AccessToken:     string(accessToken),
		},
	}
}

func FromCoreUserAuthenticToLoginResponse(shelterUser *shelter.UserAuthentic, refreshToken text.Token, accessToken pkgText.JwtToken) *AuthLoginResponse {
	var userUpdateResponse = user.FromCoreUserAuthenticToUpdateResponse(shelterUser, accessToken)
	return &AuthLoginResponse{
		UserRegisterResponse: user.UserRegisterResponse{
			UserUpdateResponse: *userUpdateResponse,
			RefreshToken:       string(refreshToken),
		},
	}
}
