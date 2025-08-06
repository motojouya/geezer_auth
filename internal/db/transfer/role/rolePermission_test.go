package role_test

import (
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/role"
	pkg "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromCoreRolePermission(t *testing.T) {
	var roleLabel, _ = pkg.NewLabel("TEST_ROLE")
	var shelterRolePermission = shelter.NewRolePermission(
		roleLabel,
		true,
		true,
		true,
		true,
		1,
	)

	var rolePermission = role.FromCoreRolePermission(shelterRolePermission)

	assert.Equal(t, string(roleLabel), rolePermission.RoleLabel)
	assert.Equal(t, rolePermission.SelfEdit, rolePermission.SelfEdit)
	assert.Equal(t, rolePermission.CompanyAccess, rolePermission.CompanyAccess)
	assert.Equal(t, rolePermission.CompanyInvite, rolePermission.CompanyInvite)
	assert.Equal(t, rolePermission.CompanyEdit, rolePermission.CompanyEdit)
	assert.Equal(t, rolePermission.Priority, rolePermission.Priority)

	t.Logf("rolePermission: %+v", rolePermission)
}

func TestToCoreRolePermission(t *testing.T) {
	var roleLabel, _ = pkg.NewLabel("TEST_ROLE")
	var rolePermission = role.RolePermission{
		RoleLabel:     string(roleLabel),
		SelfEdit:      true,
		CompanyAccess: true,
		CompanyInvite: true,
		CompanyEdit:   true,
		Priority:      1,
	}

	var shelterRolePermission, err = rolePermission.ToCoreRolePermission()
	assert.NoError(t, err)

	assert.Equal(t, roleLabel, shelterRolePermission.RoleLabel)
	assert.Equal(t, true, shelterRolePermission.SelfEdit)
	assert.Equal(t, true, shelterRolePermission.CompanyAccess)
	assert.Equal(t, true, shelterRolePermission.CompanyInvite)
	assert.Equal(t, true, shelterRolePermission.CompanyEdit)
	assert.Equal(t, uint(1), shelterRolePermission.Priority)

	t.Logf("shelterRolePermission: %+v", shelterRolePermission)
}

func TestToCoreRolePermissionError(t *testing.T) {
	var rolePermission = role.RolePermission{
		RoleLabel:     "invalid_role",
		SelfEdit:      true,
		CompanyAccess: true,
		CompanyInvite: true,
		CompanyEdit:   true,
		Priority:      1,
	}

	var _, err = rolePermission.ToCoreRolePermission()
	assert.Error(t, err, "Expected error when creating shelter role permission with invalid label")

	t.Logf("Error: %v", err)
}
