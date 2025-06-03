package role_test

import (
	"github.com/motojouya/geezer_auth/internal/core/role"
	pkg "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRolePermission(t *testing.T) {
	var roleLabel, _ = pkg.NewLabel("TEST_ROLE")
	var rolePermission = role.NewRolePermission(roleLabel, true, true, true, true, 1)

	assert.Equal(t, string(roleLabel), string(rolePermission.RoleLabel))
	assert.True(t, rolePermission.SelfEdit)
	assert.True(t, rolePermission.CompanyAccess)
	assert.True(t, rolePermission.CompanyInvite)
	assert.True(t, rolePermission.CompanyEdit)
	assert.Equal(t, uint(1), rolePermission.Priority)

	var key = role.PermissionKey(rolePermission)
	assert.Equal(t, string(roleLabel), key)

	var isRolePermission = role.PermissionIs(roleLabel)(rolePermission)
	assert.True(t, isRolePermission, "Expected rolePermission to match the label")

	var fakeLabel, _ = pkg.NewLabel("TEST_RALE")
	var isNotRolePermission = role.PermissionIs(fakeLabel)(rolePermission)
	assert.False(t, isRolePermission, "Expected rolePermission to match the label")

	t.Logf("rolePermission: %+v", rolePermission)
	t.Logf("rolePermission.RoleLabel: %s", rolePermission.RoleLabel)
	t.Logf("rolePermission.SelfEdit: %t", rolePermission.SelfEdit)
	t.Logf("rolePermission.CompanyAccess: %t", rolePermission.CompanyAccess)
	t.Logf("rolePermission.CompanyInvite: %t", rolePermission.CompanyInvite)
	t.Logf("rolePermission.CompanyEdit: %t", rolePermission.CompanyEdit)
	t.Logf("rolePermission.Priority: %d", rolePermission.Priority)
}
