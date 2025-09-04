package testUtility

import (
	entryCompanyUser "github.com/motojouya/geezer_auth/internal/entry/transfer/companyUser"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
)

type AllRoleGetterMock struct {
	FakeExecute func() ([]shelterRole.Role, error)
}

func (mock AllRoleGetterMock) Execute() ([]shelterRole.Role, error) {
	return mock.FakeExecute()
}

type RoleGetterMock struct {
	FakeExecute func(entry entryCompanyUser.RoleGetter) (*shelterRole.Role, error)
}

func (mock RoleGetterMock) Execute(entry entryCompanyUser.RoleGetter) (*shelterRole.Role, error) {
	return mock.FakeExecute(entry)
}
