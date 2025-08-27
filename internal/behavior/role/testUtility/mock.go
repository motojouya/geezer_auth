package testUtility

import (
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
)

type AllRoleGetterMock struct {
	FakeExecute func() ([]shelterRole.Role, error)
}

func (mock AllRoleGetterMock) Execute() ([]shelterRole.Role, error) {
	return mock.FakeExecute()
}

type RoleGetterMock struct {
	FakeExecute func() ([]shelterRole.Role, error)
}

func (mock RoleGetterMock) Execute() ([]shelterRole.Role, error) {
	return mock.FakeExecute()
}
