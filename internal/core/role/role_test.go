package role_test

import (
	"github.com/motojouya/geezer_auth/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateRole(t *testing.T) {
	var name = "TestRole"
	var label = "TEST_ROLE"
	var description = "Role for testing"

	var role = model.CreateRole(name, label, description)

	assert.Equal(t, name, role.Name)
	assert.Equal(t, label, role.Label)
	assert.Equal(t, description, role.Description)

	t.Logf("role: %+v", role)
	t.Logf("role.Name: %s", role.Name)
	t.Logf("role.Label: %s", role.Label)
	t.Logf("role.Description: %s", role.Description)
}

func TestNewRole(t *testing.T) {
	var roleId uint = 1
	var name = "TestRole"
	var label = "TEST_ROLE"
	var description = "Role for testing"
	var registeredDate = time.Now()

	var role = model.NewRole(roleId, name, label, description, registeredDate)

	assert.Equal(t, roleId, role.RoleId)
	assert.Equal(t, name, role.Name)
	assert.Equal(t, label, role.Label)
	assert.Equal(t, description, role.Description)
	assert.Equal(t, registeredDate, role.RegisteredDate)

	t.Logf("role: %+v", role)
	t.Logf("role.RoleId: %d", role.RoleId)
	t.Logf("role.Name: %s", role.Name)
	t.Logf("role.Label: %s", role.Label)
	t.Logf("role.Description: %s", role.Description)
	t.Logf("role.RegisteredDate: %s", role.RegisteredDate)
}
