package role_test

import (
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

	var results, err = orp.GetRole()
	if err != nil {
		t.Fatalf("Could not get roles: %s", err)
	}

	testUtility.AssertRecords(t, records, results, assertSameRole)
}

func assertSameRole(t *testing.T, expect role.Role, actual role.Role) {
	assert.Equal(t, expect.Label, actual.Label)
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Description, actual.Description)
	assert.WithinDuration(t, expect.RegisteredDate, actual.RegisteredDate, time.Second)
}
