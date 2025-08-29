package role_test

import (
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRole(t *testing.T) {
	var now = testUtility.GetNow()

	var records = []role.Role{
		//           label              , name           , description                , registeredDate
		role.NewRole("LABEL_ADMIN" /* */, "Administrator", "administrator description", now),
		role.NewRole("LABEL_MEMBER" /**/, "Member" /*  */, "member description" /*  */, now),
		role.NewRole("LABEL_STAFF" /* */, "Staff" /*   */, "staff description" /*   */, now),
	}

	testUtility.Truncate(t, orp)
	testUtility.Ready(t, orp, records)

	var results, err = orp.GetRole("LABEL_MEMBER")
	if err != nil {
		t.Fatalf("Could not get roles: %s", err)
	}

	assert.NotNil(t, results)

	assertSameRole(t, records[1], *results)
}
