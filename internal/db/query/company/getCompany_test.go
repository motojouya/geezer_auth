package company_test

import (
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/company"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetCompany(t *testing.T) {
	var now = testUtility.GetNow()

	var records = []company.Company{
		//                 persist_key,identifier,name   ,register_date
		company.NewCompany(0, "CP-TESTES", "test company", now),
		company.NewCompany(0, "CP-TASTAS", "tast company", now),
		company.NewCompany(0, "CP-TOSTOS", "tost company", now),
	}

	testUtility.Truncate(t, orp)
	var savedRecords = testUtility.Ready(t, orp, records)

	var result, err = orp.GetCompany("CP-TASTAS")
	if err != nil {
		t.Fatalf("Could not get company: %s", err)
	}

	assert.NotNil(t, result)
	assertSameCompany(t, savedRecords[1], *result)
}

func assertSameCompany(t *testing.T, expect company.Company, actual company.Company) {
	assert.Equal(t, expect.PersistKey, actual.PersistKey)
	assert.Equal(t, expect.Identifier, actual.Identifier)
	assert.Equal(t, expect.Name, actual.Name)
	assert.WithinDuration(t, expect.RegisteredDate, actual.RegisteredDate, time.Second)
}
