package company_test

import (
	"github.com/motojouya/geezer_auth/internal/db"
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"testing"
)

var orp db.ORP

func TestMain(m *testing.M) {
	testUtility.ExecuteDatabaseTest("../../../../", func(orpArg db.ORP) int {
		orp = orpArg
		return m.Run()
	})
	orp = nil // il?
}
