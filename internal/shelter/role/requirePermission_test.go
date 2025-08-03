package role_test

import (
	"github.com/motojouya/geezer_auth/internal/shelter/role"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRequirePermission(t *testing.T) {
	var selfEdit = true
	var companyAccess = false
	var companyInvite = true
	var companyEdit = false

	var requirePermission = role.NewRequirePermission(selfEdit, companyAccess, companyInvite, companyEdit)

	assert.Equal(t, selfEdit, requirePermission.SelfEdit)
	assert.Equal(t, companyAccess, requirePermission.CompanyAccess)
	assert.Equal(t, companyInvite, requirePermission.CompanyInvite)
	assert.Equal(t, companyEdit, requirePermission.CompanyEdit)

	t.Logf("requirePermission: %+v", requirePermission)
	t.Logf("requirePermission.SelfEdit: %t", requirePermission.SelfEdit)
	t.Logf("requirePermission.CompanyAccess: %t", requirePermission.CompanyAccess)
	t.Logf("requirePermission.CompanyInvite: %t", requirePermission.CompanyInvite)
	t.Logf("requirePermission.CompanyEdit: %t", requirePermission.CompanyEdit)
}
