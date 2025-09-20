package company

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type CompanyGetResponse struct {
	Company common.Company `json:"company"`
}

type CompanyTokenResponse struct {
	CompanyGetResponse
	Token string `json:"access_token"`
}

type CompanyUserResponse struct {
	Users []common.User `json:"users"`
}

func FromShelterCompany(shelterCompany shelterCompany.Company) CompanyGetResponse {
	var commonUser = common.FromShelterCompany(shelterCompany)
	return CompanyGetResponse{
		Company: commonUser,
	}
}

func FromShelterCompanyToken(shelterCompany shelterCompany.Company, accessToken pkgText.JwtToken) CompanyTokenResponse {
	var commonUser = common.FromShelterCompany(shelterCompany)
	return CompanyTokenResponse{
		CompanyGetResponse: CompanyGetResponse{
			Company: commonUser,
		},
		Token: string(accessToken),
	}
}

func FromShelterUserAuthentic(shelterUsers []shelterUser.UserAuthentic) *CompanyUserResponse {
	var ptrUsers = essence.ToPtr(shelterUsers)
	var ptrCommonUsers = essence.Map(ptrUsers, common.FromShelterUser)
	var users = essence.ToVal(ptrCommonUsers)
	return &CompanyUserResponse{
		Users: users,
	}
}
