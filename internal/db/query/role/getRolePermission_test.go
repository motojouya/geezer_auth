package role_test

import (
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRolePermission(t *testing.T) {
	var records = []role.RolePermission{
		//                      label       ,self_edit ,company_access ,company_invite ,company_edit ,priority
		role.NewRolePermission("LABEL_ADMIN", true, true, true, true, 1),
		role.NewRolePermission("LABEL_MEMBER",true, true, true, false, 2),
		role.NewRolePermission("LABEL_STAFF", true, true, false, false, 3),
	}

	testUtility.Truncate(t, orp)
	testUtility.Ready(t, orp, records)

	var results, err = orp.GetRolePermission()
	if err != nil {
		t.Fatalf("Could not get role_permission: %s", err)
	}

	testUtility.AssertRecords(t, records, results, assertSameRolePermission)
}

func assertSameRolePermission(t *testing.T, expect role.RolePermission, actual role.RolePermission) {
	assert.Equal(t, expect.RoleLabel, actual.RoleLabel)
	assert.Equal(t, expect.SelfEdit, actual.SelfEdit)
	assert.Equal(t, expect.CompanyAccess, actual.CompanyAccess)
	assert.Equal(t, expect.CompanyInvite, actual.CompanyInvite)
	assert.Equal(t, expect.CompanyEdit, actual.CompanyEdit)
	assert.Equal(t, expect.Priority, actual.Priority)
}
