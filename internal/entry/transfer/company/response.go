package company

import (
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
)

type CompanyGetResponse struct {
	Company common.Company `json:"company"`
}

type CompanyUserResponse struct {
	Users []*common.User `json:"users"`
}

func FromCoreCompany(shelterCompany shelterCompany.Company) CompanyGetResponse {
	var commonUser = common.FromCoreCompany(shelterCompany)
	return CompanyGetResponse{
		Company: commonUser,
	}
}

func FromCoreUserAuthentic(shelterUsers []*shelterUser.UserAuthentic) *CompanyUserResponse {
	var users []*common.User = essence.Map(shelterUsers, common.FromCoreUser)
	return &CompanyUserResponse{
		Users: users,
	}
}
