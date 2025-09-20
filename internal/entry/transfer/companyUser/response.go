package companyUser

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	"github.com/motojouya/geezer_auth/internal/shelter/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
)

type CompanyUserResponse struct {
	user.UserGetResponse
}

type CompanyUserTokenResponse struct {
	user.UserUpdateResponse
}

type CompanyUserInviteResponse struct {
	Token string `json:"token"`
}

type RoleGetResponse struct {
	Roles []common.Role `json:"roles"`
}

func FromShelterRoles(shelterRoles []shelterRole.Role) *RoleGetResponse {
	var roles = essence.Map(shelterRoles, common.FromShelterRole)
	return &RoleGetResponse{
		Roles: roles,
	}
}

func FromShelterUserAuthenticToGetResponse(shelterUser *shelterUser.UserAuthentic) *CompanyUserResponse {
	var commonUser = common.FromShelterUser(shelterUser)
	return &CompanyUserResponse{
		UserGetResponse: user.UserGetResponse{
			User: *commonUser,
		},
	}
}

func FromShelterUserAuthenticToTokenResponse(shelterUser *shelterUser.UserAuthentic, accessToken pkgText.JwtToken) *CompanyUserTokenResponse {
	var userResponse = user.FromShelterUserAuthenticToUpdateResponse(shelterUser, accessToken)
	return &CompanyUserTokenResponse{
		UserUpdateResponse: *userResponse,
	}
}

func FromToken(token text.Token) CompanyUserInviteResponse {
	return CompanyUserInviteResponse{
		Token: string(token),
	}
}
