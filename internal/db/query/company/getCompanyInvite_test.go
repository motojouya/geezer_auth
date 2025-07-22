package company_test

import (
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/company"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetCompanyInvite(t *testing.T) {
	var now = testUtility.GetNow()
	testUtility.Truncate(t, orp)

	var roleRecords = []role.Role{
		//           label              , name           , description                , register_date
		role.NewRole("LABEL_ADMIN" /* */, "Administrator", "administrator description", now),
		role.NewRole("LABEL_MEMBER" /**/, "Member" /*  */, "member description" /*  */, now),
		role.NewRole("LABEL_STAFF" /* */, "Staff" /*   */, "staff description" /*   */, now),
	}
	testUtility.Ready(t, orp, roleRecords)

	var companyRecords = []company.Company{
		//                 persist_key, identifier , name          , register_date
		company.NewCompany(0 /*     */, "CP-TESTES", "test company", now),
		company.NewCompany(0 /*     */, "CP-TASTAS", "tast company", now),
		company.NewCompany(0 /*     */, "CP-TOSTOS", "tost company", now),
	}
	var savedCompanyRecords = testUtility.Ready(t, orp, companyRecords)

	var companyInviteRecords = []company.CompanyInvite{
		//                       persist_key, company_persist_key              , token         , role_label         ,register_date, expire_date
		company.NewCompanyInvite(0 /*     */, savedCompanyRecords[0].PersistKey, "test_token01", "LABEL_MEMBER" /**/, now /*    */, now),
		company.NewCompanyInvite(0 /*     */, savedCompanyRecords[1].PersistKey, "test_token02", "LABEL_STAFF" /* */, now /*    */, now),
		company.NewCompanyInvite(0 /*     */, savedCompanyRecords[1].PersistKey, "test_token01", "LABEL_MEMBER" /**/, now /*    */, now),
	}
	var savedCompanyInviteRecords = testUtility.Ready(t, orp, companyInviteRecords)

	var result, err = orp.GetCompanyInvite("CP-TASTAS", "test_token01")
	if err != nil {
		t.Fatalf("Could not get company: %s", err)
	}

	assert.NotNil(t, result)
	var companyInviteExpect = company.CompanyInviteFull{
		CompanyInvite: company.CompanyInvite{
			PersistKey:        savedCompanyInviteRecords[2].PersistKey,
			CompanyPersistKey: savedCompanyRecords[1].PersistKey,
			Token:             "test_token01",
			RoleLabel:         "LABEL_MEMBER",
			RegisterDate:      now,
			ExpireDate:        now,
		},
		CompanyIdentifier:     "CP-TASTAS",
		CompanyName:           "tast company",
		CompanyRegisteredDate: now,
		RoleName:              "Member",
		RoleDescription:       "member description",
		RoleRegisteredDate:    now,
	}
	assertSameCompanyInvite(t, companyInviteExpect, *result)
}

func assertSameCompanyInvite(t *testing.T, expect company.CompanyInviteFull, actual company.CompanyInviteFull) {
	assert.Equal(t, expect.PersistKey, actual.PersistKey)
	assert.Equal(t, expect.Token, actual.Token)
	assert.WithinDuration(t, expect.RegisterDate, actual.RegisterDate, time.Second)
	assert.WithinDuration(t, expect.ExpireDate, actual.ExpireDate, time.Second)

	assert.Equal(t, expect.CompanyPersistKey, actual.CompanyPersistKey)
	assert.Equal(t, expect.CompanyIdentifier, actual.CompanyIdentifier)
	assert.Equal(t, expect.CompanyName, actual.CompanyName)
	assert.WithinDuration(t, expect.CompanyRegisteredDate, actual.CompanyRegisteredDate, time.Second)

	assert.Equal(t, expect.RoleLabel, actual.RoleLabel)
	assert.Equal(t, expect.RoleName, actual.RoleName)
	assert.Equal(t, expect.RoleDescription, actual.RoleDescription)
	assert.WithinDuration(t, expect.RoleRegisteredDate, actual.RoleRegisteredDate, time.Second)
}
