package role_test

import (
	"errors"
	behaviorRole "github.com/motojouya/geezer_auth/internal/behavior/role"
	dbRole "github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type allRoleGetterDBMock struct {
	getRole func() ([]dbRole.Role, error)
}

func (mock allRoleGetterDBMock) GetRole() ([]dbRole.Role, error) {
	return mock.getRole()
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

func getAllRoleGetterDBMock(roles []dbRole.Role) allRoleGetterDBMock {
	var getRole = func() ([]dbRole.Role, error) {
		return roles, nil
	}
	return allRoleGetterDBMock{
		getRole: getRole,
	}
}

func TestRoleGetter(t *testing.T) {
	var expectLabel01 = "ROLE_ONE"
	var role01 = getRole(expectLabel01)
	var expectLabel02 = "ROLE_TWO"
	var role02 = getRole(expectLabel02)

	var dbMock = getAllRoleGetterDBMock([]dbRole.Role{role01, role02})

	getter := behaviorRole.NewAllRoleGet(dbMock)
	result, err := getter.Execute()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result), "Expected 2 roles")
	assert.Equal(t, expectLabel01, string(result[0].Label), "Expected role label 'ROLE_ONE'")
	assert.Equal(t, expectLabel02, string(result[1].Label), "Expected role label 'ROLE_TWO'")

	t.Logf("role: %+v", result)
}

func TestRoleGetterErrGet(t *testing.T) {
	var expectLabel01 = "ROLE_ONE"
	var role01 = getRole(expectLabel01)
	var expectLabel02 = "ROLE_TWO"
	var role02 = getRole(expectLabel02)

	var dbMock = getAllRoleGetterDBMock([]dbRole.Role{role01, role02})
	dbMock.getRole = func() ([]dbRole.Role, error) {
		return []dbRole.Role{}, errors.New("database error")
	}

	getter := behaviorRole.NewAllRoleGet(dbMock)
	_, err := getter.Execute()

	assert.Error(t, err)
}

func TestRoleGetterErrTrans(t *testing.T) {
	var expectLabel01 = "ROLE_ONE"
	var role01 = getRole(expectLabel01)
	var expectLabel02 = "role_two"
	var role02 = getRole(expectLabel02)

	var dbMock = getAllRoleGetterDBMock([]dbRole.Role{role01, role02})

	getter := behaviorRole.NewAllRoleGet(dbMock)
	_, err := getter.Execute()

	assert.Error(t, err)
}
