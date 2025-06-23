package role_test

import (
	core "github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
	pkg "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromCoreRolePermission(t *testing.T) {
	var roleLabel, _ = pkg.NewLabel("TEST_ROLE")
	var coreRolePermission = core.NewRolePermission(
		roleLabel,
		true,
		true,
		true,
		true,
		1,
	)

	var rolePermission = role.FromCoreRolePermission(coreRolePermission)

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

	var coreRolePermission, err = rolePermission.ToCoreRolePermission()
	assert.NoError(t, err)

	assert.Equal(t, roleLabel, coreRolePermission.RoleLabel)
	assert.Equal(t, true, coreRolePermission.SelfEdit)
	assert.Equal(t, true, coreRolePermission.CompanyAccess)
	assert.Equal(t, true, coreRolePermission.CompanyInvite)
	assert.Equal(t, true, coreRolePermission.CompanyEdit)
	assert.Equal(t, uint(1), coreRolePermission.Priority)

	t.Logf("coreRolePermission: %+v", coreRolePermission)
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
	assert.Error(t, err, "Expected error when creating core role permission with invalid label")

	t.Logf("Error: %v", err)
}
