package role_test

import (
	coreRole "github.com/motojouya/geezer_auth/internal/core/role"
	role "github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"github.com/motojouya/geezer_auth/internal/core/text"
	pkg "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFromCoreRole(t *testing.T) {
	var name, _ = pkg.NewName("TestRole")
	var label, _ = pkg.NewLabel("TEST_ROLE")
	var description, _ = text.NewText("Role for testing")
	var registeredDate = time.Now()

	var coreRoleValue = coreRole.NewRole(name, label, description, registeredDate)
	var roleValue = role.FromCoreRole(coreRoleValue)

	assert.Equal(t, string(name), roleValue.Name)
	assert.Equal(t, string(label), roleValue.Label)
	assert.Equal(t, string(description), roleValue.Description)
	assert.Equal(t, registeredDate, roleValue.RegisteredDate)

	t.Logf("roleValue: %+v", roleValue)
}

func TestToCoreRole(t *testing.T) {
	var name = "TestRole"
	var label = "TEST_ROLE"
	var description = "Role for testing"
	var registeredDate = time.Now()

	var roleValue = role.Role{
		Name:           name,
		Label:          label,
		Description:    description,
		RegisteredDate: registeredDate,
	}

	var coreRoleValue, err = roleValue.ToCoreRole()
	assert.NoError(t, err)

	assert.Equal(t, name, string(coreRoleValue.Name))
	assert.Equal(t, label, string(coreRoleValue.Label))
	assert.Equal(t, description, string(coreRoleValue.Description))
	assert.Equal(t, registeredDate, coreRoleValue.RegisteredDate)

	t.Logf("coreRoleValue: %+v", coreRoleValue)
}

func TestToCoreRoleErrors(t *testing.T) {
	var name = "TestRole"
	var label = "invalid_role"
	var description = "Role for testing"
	var registeredDate = time.Now()

	var roleValue = role.Role{
		Name:           name,
		Label:          label,
		Description:    description,
		RegisteredDate: registeredDate,
	}

	var _, err = roleValue.ToCoreRole()
	assert.Error(t, err, "Expected error for invalid label")
}
