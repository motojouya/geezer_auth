package company

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
)

type CompanyGetResponse struct {
	Company common.Company `json:"company"`
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

func FromShelterUserAuthentic(shelterUsers []shelterUser.UserAuthentic) *CompanyUserResponse {
	var ptrUsers = essence.ToPtr(shelterUsers)
	var ptrCommonUsers = essence.Map(ptrUsers, common.FromShelterUser)
	var users = essence.ToVal(ptrCommonUsers)
	return &CompanyUserResponse{
		Users: users,
	}
}
