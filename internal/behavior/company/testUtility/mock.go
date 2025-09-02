package testUtility

import (
	entryCompany "github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	entryCompanyUser "github.com/motojouya/geezer_auth/internal/entry/transfer/companyUser"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
)

type CompanyCreatorMock struct {
	FakeExecute func(entry entryCompany.CompanyCreator) (shelterCompany.Company, error)
}

func (mock CompanyCreatorMock) Execute(entry entryCompany.CompanyCreator) (shelterCompany.Company, error) {
	return mock.FakeExecute(entry)
}

type CompanyGetterMock struct {
	FakeExecute func(entry entryCompany.CompanyGetter) (*shelterCompany.Company, error)
}

func (mock CompanyGetterMock) Execute(entry entryCompany.CompanyGetter) (*shelterCompany.Company, error) {
	return mock.FakeExecute(entry)
}

type AllUserGetterMock struct {
	FakeExecute func(entry entryCompany.CompanyGetter) ([]shelterUser.UserAuthentic, error)
}

func (mock AllUserGetterMock) Execute(entry entryCompany.CompanyGetter) ([]shelterUser.UserAuthentic, error) {
	return mock.FakeExecute(entry)
}

type UserGetterMock struct {
	FakeExecute func(entry entryCompanyUser.CompanyUserGetter, company shelterCompany.Company) (*shelterUser.UserAuthentic, error)
}

func (mock UserGetterMock) Execute(entry entryCompanyUser.CompanyUserGetter, company shelterCompany.Company) (*shelterUser.UserAuthentic, error) {
	return mock.FakeExecute(entry, company)
}

type InviteTokenCheckerMock struct {
	FakeExecute func(entry entryCompanyUser.InviteTokenGetter, company shelterCompany.Company) (shelterRole.Role, error)
}

func (mock InviteTokenCheckerMock) Execute(entry entryCompanyUser.InviteTokenGetter, company shelterCompany.Company) (shelterRole.Role, error) {
	return mock.FakeExecute(entry, company)
}

type InviteTokenIssuerMock struct {
	FakeExecute func(company shelterCompany.Company, role shelterRole.Role) (shelterText.Token, error)
}

func (mock InviteTokenIssuerMock) Execute(company shelterCompany.Company, role shelterRole.Role) (shelterText.Token, error) {
	return mock.FakeExecute(company, role)
}

type RoleAssignerMock struct {
	FakeExecute func(company shelterCompany.Company, userAuthentic *shelterUser.UserAuthentic, role shelterRole.Role) (*shelterUser.UserAuthentic, error)
}

func (mock RoleAssignerMock) Execute(company shelterCompany.Company, userAuthentic *shelterUser.UserAuthentic, role shelterRole.Role) (*shelterUser.UserAuthentic, error) {
	return mock.FakeExecute(company, userAuthentic, role)
}
