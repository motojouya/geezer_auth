package auth

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	text "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type AuthRefreshResponse struct {
	user.UserUpdateResponse
}

type AuthLoginResponse struct {
	user.UserRegisterResponse
}

func FromShelterUserAuthenticToRefreshResponse(shelterUser *shelter.UserAuthentic, accessToken pkgText.JwtToken) *AuthRefreshResponse {
	var userGetResponse = user.FromShelterUserAuthenticToGetResponse(shelterUser)
	return &AuthRefreshResponse{
		UserUpdateResponse: user.UserUpdateResponse{
			UserGetResponse: *userGetResponse,
			AccessToken:     string(accessToken),
		},
	}
}

func FromShelterUserAuthenticToLoginResponse(shelterUser *shelter.UserAuthentic, refreshToken text.Token, accessToken pkgText.JwtToken) *AuthLoginResponse {
	var userUpdateResponse = user.FromShelterUserAuthenticToUpdateResponse(shelterUser, accessToken)
	return &AuthLoginResponse{
		UserRegisterResponse: user.UserRegisterResponse{
			UserUpdateResponse: *userUpdateResponse,
			RefreshToken:       string(refreshToken),
		},
	}
}
