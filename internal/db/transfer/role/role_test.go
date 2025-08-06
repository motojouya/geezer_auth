package role_test

import (
	role "github.com/motojouya/geezer_auth/internal/db/transfer/role"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	"github.com/motojouya/geezer_auth/internal/shelter/text"
	pkg "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFromCoreRole(t *testing.T) {
	var name, _ = pkg.NewName("TestRole")
	var label, _ = pkg.NewLabel("TEST_ROLE")
	var description, _ = text.NewText("Role for testing")
	var registeredDate = time.Now()

	var shelterRoleValue = shelterRole.NewRole(name, label, description, registeredDate)
	var roleValue = role.FromCoreRole(shelterRoleValue)

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

	var shelterRoleValue, err = roleValue.ToCoreRole()
	assert.NoError(t, err)

	assert.Equal(t, name, string(shelterRoleValue.Name))
	assert.Equal(t, label, string(shelterRoleValue.Label))
	assert.Equal(t, description, string(shelterRoleValue.Description))
	assert.Equal(t, registeredDate, shelterRoleValue.RegisteredDate)

	t.Logf("shelterRoleValue: %+v", shelterRoleValue)
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
