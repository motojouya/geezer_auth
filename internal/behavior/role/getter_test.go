package role_test

import (
	"errors"
	behaviorRole "github.com/motojouya/geezer_auth/internal/behavior/role"
	dbRole "github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type roleGetterDBMock struct {
	getRole func(label string) (*dbRole.Role, error)
}

func (mock roleGetterDBMock) GetRole(label string) (*dbRole.Role, error) {
	return mock.getRole(label)
}

func getRole(expectLabel string) dbRole.Role {
	var name = "TestRole"
	var label = expectLabel
	var description = "Role for testing"
	var registeredDate = time.Now()

	return dbRole.Role{
		Name:           name,
		Label:          label,
		Description:    description,
		RegisteredDate: registeredDate,
	}
}

func getRoleGetterDBMock(role *dbRole.Role) roleGetterDBMock {
	var getRole = func() (*dbRole.Role, error) {
		return role, nil
	}
	return roleGetterDBMock{
		getRole: getRole,
	}
}

func TestRoleGetter(t *testing.T) {
	var expectLabel01 = "ROLE_ONE"
	var role01 = getRole(expectLabel01)

	var dbMock = getRoleGetterDBMock(&role01)

	getter := behaviorRole.NewRoleGet(dbMock)
	result, err := getter.Execute()

	assert.NoError(t, err)
	assert.NotNil(t, result, "Expected non-nil result")
	assert.Equal(t, expectLabel01, string(result.Label), "Expected role label 'ROLE_ONE'")

	t.Logf("role: %+v", result)
}

func TestRoleGetterNil(t *testing.T) {
	var dbMock = getRoleGetterDBMock(nil)

	getter := behaviorRole.NewRoleGet(dbMock)
	result, err := getter.Execute()

	assert.NoError(t, err)
	assert.Nil(t, result, "Expected non-nil result")
}

func TestRoleGetterErrGet(t *testing.T) {
	var expectLabel01 = "ROLE_ONE"
	var role01 = getRole(expectLabel01)

	var dbMock = getRoleGetterDBMock(&role01)
	dbMock.getRole = func() (*dbRole.Role, error) {
		return nil, errors.New("database error")
	}

	getter := behaviorRole.NewRoleGet(dbMock)
	_, err := getter.Execute()

	assert.Error(t, err)
}

func TestRoleGetterErrTrans(t *testing.T) {
	var expectLabel01 = "role_one"
	var role01 = getRole(expectLabel01)

	var dbMock = getRoleGetterDBMock(&role01)

	getter := behaviorRole.NewRoleGet(dbMock)
	_, err := getter.Execute()

	assert.Error(t, err)
}
