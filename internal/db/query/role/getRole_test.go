package role_test

import (
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"github.com/stretchr/testify/assert"
	"testing"
)

var now = testUtility.GetNow()

var records = []role.Role{
	// カラムの位置をあわせてtableっぽくしたかったがformatterが邪魔
	//           label        , name           , description                , registeredDate
	role.NewRole("LABEL_ADMIN", "Administrator", "administrator description", now),
	role.NewRole("LABEL_MEMBER", "Member", "member description", now),
}

func TestGetRole(t *testing.T) {
	testUtility.Truncate(t, orp)
	testUtility.Ready(t, orp, records...)

	var results, err = orp.GetRole()
	if err != nil {
		t.Fatalf("Could not get roles: %s", err)
	}

	assert.ElementsMatch(t, records, results)
}
